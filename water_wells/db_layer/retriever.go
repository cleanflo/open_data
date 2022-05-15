package db_layer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/schema"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var gormLogger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags),
	logger.Config{
		SlowThreshold: time.Second,
		LogLevel:      logger.Info,
		Colorful:      true,
	})

var gormDB map[string]*gorm.DB

func (R *Retreiver) Session() (*gorm.DB, error) {
	if db, ok := gormDB[R.DBName]; ok && db != nil {
		return db, nil
	}

	// get DB connection
	query := url.Values{}
	query.Add("database", R.DBName)

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(R.user, R.password),
		Host:     fmt.Sprintf("%s:%d", R.host, R.port),
		RawQuery: query.Encode(),
	}

	dsn := u.String()
	return gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger:      gormLogger,
		DryRun:      false,
		QueryFields: true,
	})
}

type Retreiver struct {
	DBName      string                               // name of the DB
	Table       string                               // name of table in DB
	Columns     []string                             // columns to return
	Joins       []JoinStatement                      // joins required for query
	childJoins  []JoinStatement                      // joins required by query options
	Query       ClientRequestData                    // options for query params
	Data        interface{}                          // the struct of the columns listed above
	Coordinates func(data interface{}) []Coordinates // function to minify coordinate size

	ResponseChunks map[int]int // map of query size to optimum page count for this DB

	user     string
	password string
	host     string
	port     int
}

func (R *Retreiver) User(u string) *Retreiver {
	R.user = u
	return R
}

func (R *Retreiver) Password(p string) *Retreiver {
	R.password = p
	return R
}

func (R *Retreiver) Host(h string) *Retreiver {
	R.host = h
	return R
}

func (R *Retreiver) Port(p int) *Retreiver {
	R.port = p
	return R
}

var decoder = schema.NewDecoder()

func (R *Retreiver) Handler(w http.ResponseWriter, r *http.Request) {
	// move to middleware
	requestStart := time.Now()

	// Parse the request form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse request: %s", err), http.StatusBadRequest)
		return
	}

	// decode request data
	requestData := ClientRequestData{}
	// r.PostForm is a map of our POST form values
	err = decoder.Decode(&requestData, r.Form)
	if err != nil {
		// Handle error
		http.Error(w, fmt.Sprintf("Failed to decode request: %s", err), http.StatusBadRequest)
		return
	}

	queryData := R.Query.MergeRequest(requestData)
	cacheHit := false
	// check if request provides page#
	if requestData.Page != 0 {
		// check for response in cache and return
		if cacheHit {
			return
		}
	}

	db, err := R.Session()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to connect database %s", err), http.StatusInternalServerError)
		return
	}

	// build query
	tx := db.Table(R.Table)
	R.BuildQuery(queryData, tx)

	// If any global joins or joins added by the build query step
	if len(R.Joins) > 0 || len(R.childJoins) > 0 {
		builtJoins := append(R.Joins, R.childJoins...)

		fmt.Println("Adding Joins to query:", len(R.Joins))
		joinKeys := []string{}
		joinList := map[string]JoinStatement{}
		for i, join := range builtJoins {
			if j, ok := joinList[join.Right.Name]; ok {
				if j.CorrelationName() != join.CorrelationName() {
					joinKeys = append(joinKeys, fmt.Sprintf("%d*%s", i, join.CorrelationName()))
					joinList[join.CorrelationName()] = join
				}
			} else {
				joinKeys = append(joinKeys, fmt.Sprintf("%d*%s", i, join.Right.Name))
				joinList[join.Right.Name] = join
			}
		}

		sort.Strings(joinKeys)
		for i, k := range joinKeys {
			strings.Index(k, "*")
			cName := strings.Split(k, "*")[1]

			join := joinList[cName]
			tx = tx.Joins(join.Statement(cName))
			if i == 0 {
				tx = tx.Select(join.Select)
				if join.GroupBy != "" {
					tx = tx.Group(join.GroupBy)
				}
				if join.OrderBy != "" {
					tx = tx.Order(fmt.Sprintf("%s DESC", join.OrderBy))
				}
			}
		}
	}

	if requestData.Page == 0 {
		queryData.Page = 1
	} else {
		queryData.Page = requestData.Page
	}

	// check total, count in request
	// if it has, skip getting count
	queryData.Total, queryData.Chunk = requestData.Total, requestData.Chunk
	if queryData.Total == 0 || queryData.Chunk == 0 {
		// if it doesn't have, check for count in cache
		// if requestData.Total > 1000 {
		// 	cacheHit = true
		// }
	}

	// if it doesn't have count in cache and not in request
	// get total records that match the query
	if !cacheHit && queryData.Total == 0 {
		tx.Count(&queryData.Total)

		// if total == 0 then return and cache

		// save count to cache
	}

	// calculate ideal chunk size
	// we have no cache hit, a total with content, no count
	if queryData.Total > 0 {
		// calculate ideal chunk size
		if len(R.ResponseChunks) > 0 {
			// check list of chunks on retreiver to find ideal page count for query size
			largestPageCount := 0
			largestQuerySize := 0
		ResponseChunkLoop:
			for rQuerySize, rPageCount := range R.ResponseChunks {
				if rQuerySize > largestQuerySize {
					largestQuerySize = rQuerySize
				}
				if rPageCount > largestPageCount {
					largestPageCount = rPageCount
				}
				// if total is less than max query size, and pagecount is greater than this map element
				if int(queryData.Total) <= rQuerySize && (queryData.PageCount == 0 || queryData.PageCount > rPageCount) {
					queryData.PageCount = rPageCount
					continue ResponseChunkLoop
				}
			}
			if queryData.PageCount == 0 {
				queryData.PageCount = int(((queryData.Total / int64(largestQuerySize)) + 1) * int64(largestPageCount))
			}
		} else {
			// if user hasn't defined chunk sizes then just use some default sizes
			switch true {
			case queryData.Total <= 10000:
				queryData.PageCount = 1
			case queryData.Total <= 50000:
				queryData.PageCount = 2
			case queryData.Total <= 100000:
				queryData.PageCount = 5
			case queryData.Total <= 250000:
				queryData.PageCount = 10
			default:
				queryData.PageCount = 25
			}
		}

		// save the chunk size
		queryData.Chunk = queryData.Total / int64(queryData.PageCount)

		// find size of last chunk
		lastChunk := queryData.Total % queryData.Chunk
		// if last chunk is less than 25% of the chunksize
		if lastChunk > 0 && lastChunk < (queryData.Chunk/4) {
			// then spread the last page amongst other pages
			queryData.Chunk += lastChunk / int64(queryData.PageCount)
		} else if lastChunk > 0 {
			queryData.PageCount += 1
		}
	}

	// run query for chunk size
	g := tx.Session(&gorm.Session{
		Logger: gormLogger,
	})

	fmt.Printf("Calculated total(%d)\t chunk (%d)\t page(%d)\t\n", queryData.Total, queryData.Chunk, queryData.Page)

	// retrieve and return requested data
	// gofunc retreive next page of data and store in cache
	g = g.Limit(int(queryData.Chunk)).Offset(int(queryData.Chunk) * (queryData.Page - 1))
	g.Find(R.Data)

	// run built query through SQL
	dry := g.Session(&gorm.Session{
		Logger: gormLogger,
		DryRun: true,
	})
	stmt := dry.Find(R.Data).Statement

	stmtString := stmt.SQL.String()
	for i, a := range stmt.Vars {
		paramPosition := fmt.Sprintf("@p%d", i+1)
		stmtString = strings.Replace(stmtString, paramPosition, fmt.Sprint(a), 1)
	}
	queryData.Statement = stmtString

	// retreive next page of data in a go func
	// go func() {
	// 	gn := tx.Session(&gorm.Session{
	// 		Logger: gormLogger,
	// 	})
	// 	gn = gn.Limit(int(pageChunkSize)).Offset(int(pageChunkSize) * page)
	// 	gn.Find(&R.Data)
	// }()

	pointList := R.ProcessData()

	resp := ClientResponse{
		DefaultResponse: queryData.DefaultResponse,
		Data:            pointList,
	}

	// save result to cache

	// return result as json
	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal response %s", err), http.StatusInternalServerError)
		return
	}
	fmt.Printf("JSON size: %d\n", len(b))

	w.Header().Set("X-Content-Type-Options", "nosniff")
	if _, err = w.Write(b); err != nil {
		http.Error(w, fmt.Sprintf("failed to write response %s", err), http.StatusInternalServerError)
		return
	}

	d := time.Since(requestStart)
	fmt.Printf("Request time: %s\n", d)
}

func (R *Retreiver) ProcessData() (c []Coordinates) {
	return R.Coordinates(R.Data)
}
