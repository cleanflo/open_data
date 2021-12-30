package db_layer

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type JoinDiffType = string

const (
	JOIN       JoinDiffType = "JOIN"
	INNER_JOIN JoinDiffType = "INNER JOIN"
	LEFT_JOIN  JoinDiffType = "LEFT JOIN"
	RIGHT_JOIN JoinDiffType = "RIGHT JOIN"
)

type JoinTable struct {
	Name string
	Key  string
}

type JoinStatement struct {
	Diff    JoinDiffType
	Left    JoinTable
	Right   JoinTable
	GroupBy string
	OrderBy string
	Select  string
}

func (j JoinStatement) CorrelationName() string {
	return fmt.Sprintf("%s%s%s%s",
		strings.ToLower(j.Left.Name),
		strings.ToTitle(j.Left.Key),
		strings.ToTitle(j.Right.Name),
		strings.ToTitle(j.Right.Key),
	)
}

func (j JoinStatement) Statement(name string) string {
	cName := ""
	if name != "" && name != j.Right.Name {
		cName = fmt.Sprintf("AS %s ", name)
	}
	return fmt.Sprintf("%s %s %s ON %s.%s = %s.%s",
		j.Diff, j.Right.Name, cName, name, j.Right.Key, j.Left.Name, j.Left.Key)
}

// Optioner is the interface definition that decodes the request, builds the query, and encodes the respose for each query paramater
// Optioner acts as the layer between a get request parameter and the DB request builder
type Optioner interface {
	Request() RequestOption
	column() string
	required() map[string]interface{}
	joins() []JoinStatement
}

// Null Option is an empty struct that fulfills the Optioner interface, it does nothing
type NullOption struct {
	// column string
}

func (n NullOption) request() RequestOption {
	return nil
}

func (n NullOption) column() string {
	return ""
}

func (n NullOption) required() map[string]interface{} {
	return nil
}

func (n NullOption) joins() []JoinStatement {
	return []JoinStatement{}
}

// TimeOption is a struct that fulfills the Optioner interface, it decodes start and end time and
// builds a query with a defined layout.
type TimeOption struct {
	Column   string                 `json:"-" schema:"-"`
	Required map[string]interface{} `json:"-" schema:"-"`
	Joins    []JoinStatement        `json:"-" schema:"-"`
	Layout   string                 `json:"-" schema:"-"`
	Start    time.Time              `json:"-" schema:"start"`
	End      time.Time              `json:"-" schema:"end"`
}

func (t TimeOption) required() map[string]interface{} {
	return t.Required
}

func (t TimeOption) joins() []JoinStatement {
	return t.Joins
}

func (t TimeOption) column() string {
	return t.Column
}

func (t TimeOption) Merge(src TimeOption) TimeOption {
	return TimeOption{
		Column:   t.Column,   // receiver
		Required: t.Required, // receiver
		Joins:    t.Joins,    // receiver
		Layout:   t.Layout,   // receiver
		Start:    src.Start,  // src
		End:      src.End,    // src
	}
}

// Query builds a DB query
func (t TimeOption) Request() RequestOption {
	start := t.Start.Format(t.Layout)
	if t.Start.IsZero() {
		start = ""
	}
	end := t.End.Format(t.Layout)
	if t.End.IsZero() {
		end = ""
	}
	return RequestRange{
		Start: start,
		End:   end,
	}
}

// UnmarshalText fulfills the interface for github.com/gorilla/schema
func (t *TimeOption) UnmarshalText(text []byte) (err error) {
	fmt.Println("TimeOption got text: ", string(text))
	return nil
}

// UnmarshalJSON fulfills the interface for encoding/json
func (t *TimeOption) MarshalJSON() (b []byte, err error) {
	return b, nil
}

// NumberOption is a struct that fulfills the Optioner interface, it decodes start and end number and
// builds a query with a BETWEEN.
type NumberOption struct {
	Column   string                 `json:"-" schema:"-"`
	Required map[string]interface{} `json:"-" schema:"-"`
	Joins    []JoinStatement        `json:"-" schema:"-"`
	Start    int                    `json:"-" schema:"start"`
	End      int                    `json:"-" schema:"end"`
}

func (n NumberOption) required() map[string]interface{} {
	return n.Required
}

func (n NumberOption) joins() []JoinStatement {
	return n.Joins
}

func (n NumberOption) column() string {
	return n.Column
}

func (n NumberOption) Merge(src NumberOption) NumberOption {
	return NumberOption{
		Column:   n.Column,
		Required: n.Required,
		Joins:    n.Joins,
		Start:    src.Start,
		End:      src.End,
	}
}

func (n NumberOption) Request() RequestOption {
	return RequestRange{n.Start, n.End, UNSET}
}

// UnmarshalText fulfills the interface for github.com/gorilla/schema
func (n *NumberOption) UnmarshalText(text []byte) (err error) {
	fmt.Println(string(text))
	v := strings.Split(string(text), ":")
	if len(v) >= 1 {
		i, err := strconv.Atoi(v[0])
		if err != nil {
			return fmt.Errorf("failed to decode number for option %s: %s", n.Column, v[0])
		}
		n.Start = i
	}
	if len(v) >= 2 {
		i, err := strconv.Atoi(v[0])
		if err != nil {
			return fmt.Errorf("failed to decode number for option %s: %s", n.Column, v[1])
		}
		n.End = i
	}
	return nil
}

// UnmarshalJSON fulfills the interface for encoding/json
func (n NumberOption) MarshalJSON() (b []byte, err error) {
	return b, nil
}

// NumberListOption is a struct that fulfills the Optioner interface, it decodes list of values and
// builds a query with a IN from the list of DB specific equivalents.
type NumberListOption struct {
	Column   string                   `json:"-" schema:"-"`
	Required map[string]interface{}   `json:"-" schema:"-"`
	Joins    []JoinStatement          `json:"-" schema:"-"`
	List     []string                 `json:"-" schema:"-"` // values passed from client for filtering
	Items    map[string][]interface{} `json:"-" schema:"-"` // map of possible client filters for SQL
	Multiple bool                     `json:"-" schema:"-"` // allow multiple filters
}

func (n NumberListOption) required() map[string]interface{} {
	return n.Required
}

func (n NumberListOption) joins() []JoinStatement {
	return n.Joins
}

func (n NumberListOption) column() string {
	return n.Column
}

// UnmarshalText fulfills the interface for github.com/gorilla/schema
func (n *NumberListOption) UnmarshalText(text []byte) (err error) {
	s := strings.Split(string(text), ",")
	n.List = s
	return nil
}

// UnmarshalJSON fulfills the interface for encoding/json
func (n *NumberListOption) MarshalJSON() (b []byte, err error) {
	return b, nil
}

func (n NumberListOption) toInterfaceSlice() map[string][]interface{} {
	m := make(map[string][]interface{})
	for k, v := range n.Items {
		i := []interface{}{}
		i = append(i, v...)
		m[k] = i
	}
	return m
}

// Query retreive lists' of integers that map to strings provided
func (n NumberListOption) Request() RequestOption {
	r := RequestList{}
ITEMS:
	for k, v := range n.Items {
		for i, a := range n.List {
			if k == a {
				// copy items to list
				r.list = append(r.list, v...)
				// remove key from origin list
				n.List = append(n.List[:i], n.List[i+1:]...)
				if !n.Multiple {
					break ITEMS
				}
				continue ITEMS
			}
		}
		// copy items to not list
		r.not = append(r.not, v...)
	}

	return r
}

func (n NumberListOption) Merge(src NumberListOption) (v NumberListOption) {
	return NumberListOption{
		Column:   n.Column,
		Required: n.Required,
		Joins:    n.Joins,
		Multiple: n.Multiple,
		Items:    n.Items,
		List: func() []string {
			s := make([]string, len(src.List))
			return append(s, src.List...)
		}(),
	}
}

type QueryType = int

const (
	UNSET QueryType = 0 + iota
	EQUAL
	NOT_EQUAL
	LESS_THAN
	GREATER_THAN
	IN_ARRAY
	NOT_ARRAY
	BETWEEN
)

type Query struct {
	query     string
	value     interface{}
	operation QueryType
}

type RequestOption interface {
	Query(column string) (Query, error)
}

type RequestList struct {
	list    []interface{}
	not     []interface{}
	hasNull bool
}

func (r RequestList) Query(column string) (q Query, err error) {
	listLen := len(r.list)
	notLen := len(r.not)
	itemList := []interface{}{}
	switch true {
	case listLen == 0:
		return q, nil
	case listLen == 1 && notLen == 0:
		// op EQUAL
		// {value} == list[0]
		q.operation = EQUAL
		q.query = fmt.Sprintf("%s = ?", column)
		q.value = r.list[0]
		itemList = append(itemList, r.list[0])
	case listLen == 0 && notLen == 1:
		// op NOT_EQUAL
		// {value} != not[0]
		q.operation = NOT_EQUAL
		q.query = fmt.Sprintf("%s = ?", column)
		q.value = r.not[0]
		itemList = append(itemList, r.not[0])
	case listLen > 1 && notLen <= 0:
		// op IN
		// {value} IN list[]
		q.operation = IN_ARRAY
		q.query = fmt.Sprintf("%s IN ?", column)
		q.value = r.list
		itemList = r.list
	case listLen <= 0 && notLen > 1:
		// op NOT IN
		// {value} NOT IN not[]
		q.operation = NOT_ARRAY
		q.query = fmt.Sprintf("%s NOT IN ?", column)
		q.value = r.not
		itemList = r.not
	case listLen > 1 && notLen > 1 && listLen < notLen:
		// op IN
		// {value} IN list[]
		q.operation = IN_ARRAY
		q.query = fmt.Sprintf("%s IN ?", column)
		q.value = r.list
		itemList = r.list
	case listLen > 1 && notLen > 1 && listLen > notLen:
		// op NOT IN
		// {value} NOT IN not[]
		q.operation = NOT_ARRAY
		q.query = fmt.Sprintf("%s NOT IN ?", column)
		q.value = r.not
		itemList = r.not
	}

	// find null items listed in the selected parameter options and remove from the list
findNull:
	for i, a := range itemList {
		if s, ok := a.(string); a == nil || ok && (s == "NULL" || s == "null") {
			r.hasNull = true
			itemList = append(itemList[:i], itemList[i+1:]...)
			// itemList[i] = itemList[len(itemList)-1]
			// itemList = itemList[:len(itemList)-1]
			break findNull
		}
	}

	// add null checks appropriately to the query
	itemListLen := len(itemList)
	switch true {
	case itemListLen == 0 && !r.hasNull:
		// no values at all
	case itemListLen == 0 && r.hasNull:
		// only value is null
		q.query = fmt.Sprintf(" %s IS NULL", column)
	case itemListLen > 0 && !r.hasNull:
		// values but no null
	case itemListLen > 0 && r.hasNull:
		// values and null
		n := ""
		if listLen > notLen {
			n = "NOT"
		}
		q.query += fmt.Sprintf(" OR %s IS %s NULL", column, n)
	}

	return q, nil
}

type RequestRange struct {
	Start     interface{}
	End       interface{}
	queryType QueryType
}

func (r RequestRange) Query(column string) (q Query, err error) {
	var startNil, endNil bool
	for i, a := range []interface{}{r.Start, r.End} {
		t := false
		switch v := a.(type) {
		case string:
			if v == "" {
				t = true
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			if v == 0 {
				t = true
			}
		case float32, float64:
			if v == 0.0 {
				t = true
			}
		case nil:
			t = true
		}

		switch i {
		case 0:
			startNil = t
		case 1:
			endNil = t
		}
	}

	if startNil && endNil {
		return q, fmt.Errorf("nil start and end")
	}

	if endNil {
		// return fmt.Sprintf("%s >= ?", column), r.Start, nil
		q.operation = GREATER_THAN
		q.query = fmt.Sprintf("%s >= ?", column)
		q.value = r.Start
		return q, nil
	}
	q.operation = BETWEEN
	q.query = fmt.Sprintf("%s BETWEEN ? AND ?", column)
	q.value = []interface{}{r.Start, r.End}
	return q, nil
}
