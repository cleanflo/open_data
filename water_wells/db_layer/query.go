package db_layer

import "fmt"

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
	case (listLen > 1 && notLen <= 0) || (listLen > 1 && notLen > 1 && listLen < notLen):
		// op IN
		// {value} IN list[]
		q.operation = IN_ARRAY
		q.query = fmt.Sprintf("%s IN ?", column)
		q.value = r.list
		itemList = r.list
	case (listLen <= 0 && notLen > 1) || (listLen > 1 && notLen > 1 && listLen > notLen):
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

	if startNil {
		q.operation = LESS_THAN
		q.query = fmt.Sprintf("%s <= ?", column)
		q.value = r.End
		return q, nil
	}

	q.operation = BETWEEN
	q.query = fmt.Sprintf("%s BETWEEN ? AND ?", column)
	q.value = []interface{}{r.Start, r.End}
	return q, nil
}
