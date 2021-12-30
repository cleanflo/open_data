package db_layer

import (
	"reflect"

	"gorm.io/gorm"
)

type Coordinates []float32

func (c *Coordinates) Add(lat, lng float32) {
	*c = []float32{lat, lng}
}

type DefaultResponse struct {
	Statement string `json:"statement" schema:"-"`
	Total     int64  `json:"total" schema:"total"`
	Chunk     int64  `json:"chunk" schema:"chunk"`
	Page      int    `json:"page" schema:"page"`
	PageCount int    `json:"pageCount" schema:"-"`
}

type ErrorResponse struct {
	DefaultResponse
	Error error `json:"error"`
}

type ClientResponse struct {
	DefaultResponse
	Data []Coordinates `json:"data" schema:"-"`
}

type ClientRequestData struct {
	ClientResponse

	Completed TimeOption       `json:"-" schema:"completed"`
	Abandoned TimeOption       `json:"-" schema:"abandoned"`
	Status    NumberListOption `json:"-" schema:"status"`
	Use       NumberListOption `json:"-" schema:"use"`
	Colour    NumberListOption `json:"-" schema:"colour"`
	Taste     NumberListOption `json:"-" schema:"taste"`
	Odour     NumberListOption `json:"-" schema:"odour"`
	Rate      NumberOption     `json:"-" schema:"rate"`
	Depth     NumberOption     `json:"-" schema:"depth"`
	Bedrock   NumberOption     `json:"-" schema:"bedrock"`
}

func (c ClientRequestData) MergeRequest(src ClientRequestData) ClientRequestData {
	c.Completed = c.Completed.Merge(src.Completed)
	c.Abandoned = c.Abandoned.Merge(src.Abandoned)
	c.Status = c.Status.Merge(src.Status)
	c.Use = c.Use.Merge(src.Use)
	c.Colour = c.Colour.Merge(src.Colour)
	c.Taste = c.Taste.Merge(src.Taste)
	c.Odour = c.Odour.Merge(src.Odour)
	c.Rate = c.Rate.Merge(src.Rate)
	c.Depth = c.Depth.Merge(src.Depth)
	c.Bedrock = c.Bedrock.Merge(src.Bedrock)
	return c
}

func (R *Retreiver) BuildQuery(q ClientRequestData, db *gorm.DB) *gorm.DB {
	val := reflect.ValueOf(q)
	if val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	R.childJoins = []JoinStatement{}
	for i := 0; i < val.NumField(); i += 1 {
		field := val.Field(i)
		// fieldType := field.Type()

		if opt, ok := field.Interface().(Optioner); ok {
			// execute build query
			req := opt.Request()
			q, err := req.Query(opt.column())
			if err != nil {
				// fmt.Printf("Failed to build query: %s\n", err)
				continue
			}

			var tx *gorm.DB
			switch q.operation {
			case EQUAL:
				tx = db.Where(q.query, q.value)
			case NOT_EQUAL:
				tx = db.Not(q.query, q.value)
			case IN_ARRAY:
				tx = db.Where(q.query, q.value)
			case NOT_ARRAY:
				tx = db.Where(q.query, q.value)
				// tx = db.Not(map[string]interface{}{q.query: q.value})
			case GREATER_THAN:
				tx = db.Where(q.query, q.value)
			case BETWEEN:
				if v, ok := q.value.([]interface{}); ok {
					tx = db.Where(q.query, v...)
				}
			}

			if tx != nil {
				if required := opt.required(); required != nil {
					tx = db.Where(required)
				}
				if joins := opt.joins(); len(joins) > 0 {
					// copy joins to receiver
					R.childJoins = append(R.childJoins, joins...)
				}
				db = tx
			}
		}
	}
	return db
}
