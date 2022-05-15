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

func capitalize(s string) string {
	if len(s) < 2 {
		return s
	}
	return strings.ToUpper(s[0:1]) + strings.ToLower(s[1:])
}

func (j JoinStatement) CorrelationName() string {
	return fmt.Sprintf("%s%s%s%s",
		strings.ToLower(j.Left.Name),
		capitalize(j.Left.Key),
		capitalize(j.Right.Name),
		capitalize(j.Right.Key),
	)
}

func (j JoinStatement) Statement(name string) string {
	cName := ""
	if name == "" {
		name = j.Right.Name
	}

	if name != j.Right.Name {
		cName = fmt.Sprintf("AS %s", name)
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
// type NullOption struct {
// 	// column string
// }

// func (n NullOption) request() RequestOption {
// 	return nil
// }

// func (n NullOption) column() string {
// 	return ""
// }

// func (n NullOption) required() map[string]interface{} {
// 	return nil
// }

// func (n NullOption) joins() []JoinStatement {
// 	return []JoinStatement{}
// }

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
	v := strings.Split(string(text), ":")
	if len(v) >= 1 {
		i, err := strconv.Atoi(v[0])
		if err != nil {
			return fmt.Errorf("failed to decode number for option %s: %s", n.Column, v[0])
		}
		n.Start = i
	}
	if len(v) >= 2 {
		i, err := strconv.Atoi(v[1])
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

// func (n NumberListOption) toInterfaceSlice() map[string][]interface{} {
// 	m := make(map[string][]interface{})
// 	for k, v := range n.Items {
// 		i := []interface{}{}
// 		i = append(i, v...)
// 		m[k] = i
// 	}
// 	return m
// }

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
			return append([]string{}, src.List...)
		}(),
	}
}
