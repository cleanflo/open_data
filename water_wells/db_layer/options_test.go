package db_layer

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testJoinStatement = JoinStatement{
	Diff:    INNER_JOIN,
	Left:    JoinTable{"lname", "lkey"},
	Right:   JoinTable{"rname", "rkey"},
	GroupBy: "gtest",
	OrderBy: "otest",
	Select:  "seltest",
}

func TestJoinStatement(t *testing.T) {
	j := testJoinStatement
	t.Run("CorrelationName", func(t *testing.T) {
		assert.Equalf(t, "lnameLkeyRnameRkey", j.CorrelationName(), "CorrelationName failed")
	})

	t.Run("Statement", func(t *testing.T) {
		s := j.Statement("")
		assert.Equalf(t, "INNER JOIN rname  ON rname.rkey = lname.lkey", s, "Statement failed")

		s = j.Statement(j.Right.Name)
		assert.Equalf(t, "INNER JOIN rname  ON rname.rkey = lname.lkey", s, "Statement failed")

		s = j.Statement("colname")
		assert.Equalf(t, "INNER JOIN rname AS colname ON colname.rkey = lname.lkey", s, "Statement failed")
	})
}

var testTimeOption = TimeOption{
	Column: "table.column",
	Layout: "2006-01-02 3:04:00.000 PM",
	Joins: []JoinStatement{
		{
			Diff:    INNER_JOIN,
			Left:    JoinTable{"lkey", "lcol"},
			Right:   JoinTable{"rkey", "rcol"},
			GroupBy: "leky.lcol",
			OrderBy: "MAX(id)",
			Select:  "MAX(id) as last, MIN(id) as first",
		},
	},
}

func TestTimeOption(t *testing.T) {
	t.Run("TimeOption", func(t *testing.T) {
		to := testTimeOption

		req := to.Request()
		assert.IsTypef(t, RequestRange{}, req, "TimeOption failed on zero-time Request: %T", req)

		rr, ok := req.(RequestRange)
		assert.Truef(t, ok, "TimeOption failed on RequestRange cast: %T", req)

		assert.Equal(t, "", rr.Start, "TimeOption failed on zero-time RequestRange.Start")
		assert.Equal(t, "", rr.End, "TimeOption failed on zero-time RequestRange.End")

		start, end := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC), time.Now()
		tm := to.Merge(TimeOption{
			Start: start,
			End:   end,
		})

		assert.Truef(t, tm.Start.Equal(start), "TimeOption.Merge failed on start: %s || %s", to.Start, start)
		assert.Truef(t, tm.End.Equal(end), "TimeOption.Merge failed on end: %s || %s", to.End, end)

		req = tm.Request()
		assert.IsTypef(t, RequestRange{}, req, "TimeOption failed on non-zero-time Request: %T", req)

		rr, ok = req.(RequestRange)
		assert.Truef(t, ok, "TimeOption failed on RequestRange cast: %T", req)

		assert.Equal(t, tm.Start.Format(tm.Layout), rr.Start, "TimeOption failed on format time RequestRange.Start")
		assert.Equal(t, tm.End.Format(tm.Layout), rr.End, "TimeOption failed on format time RequestRange.End")

	})
}

var testNumberOption = NumberOption{
	Column: "table.column",
	Joins: []JoinStatement{
		{
			Diff:    INNER_JOIN,
			Left:    JoinTable{"lkey", "lcol"},
			Right:   JoinTable{"rkey", "rcol"},
			GroupBy: "leky.lcol",
			OrderBy: "MAX(id)",
			Select:  "MAX(id) as last, MIN(id) as first",
		},
	},
}

func TestNumberOption(t *testing.T) {
	t.Run("NumberOption", func(t *testing.T) {
		no := testNumberOption

		req := no.Request()
		assert.IsTypef(t, RequestRange{}, req, "NumberOption failed on zero-value Request: %T", req)

		rr, ok := req.(RequestRange)
		assert.Truef(t, ok, "NumberOption failed on RequestRange cast: %T", req)

		assert.Equal(t, 0, rr.Start, "NumberOption failed on zero-value RequestRange.Start")
		assert.Equal(t, 0, rr.End, "NumberOption failed on zero-value RequestRange.End")

		start, end := 50, 500
		nm := no.Merge(NumberOption{
			Start: start,
			End:   end,
		})

		assert.Equalf(t, start, nm.Start, "NumberOption.Merge failed on start: %d || %d", no.Start, start)
		assert.Equalf(t, end, nm.End, "NumberOption.Merge failed on end: %d || %d", no.End, end)

		req = nm.Request()
		assert.IsTypef(t, RequestRange{}, req, "NumberOption failed on non-zero-value Request: %T", req)

		rr, ok = req.(RequestRange)
		assert.Truef(t, ok, "NumberOption failed on RequestRange cast: %T", req)

		assert.Equal(t, nm.Start, rr.Start, "NumberOption failed on RequestRange.Start")
		assert.Equal(t, nm.End, rr.End, "NumberOption failed on RequestRange.End")

		err := no.UnmarshalText([]byte(fmt.Sprintf("%d:%d", start+5, end+5)))
		assert.Nil(t, err, "NumberOption failed on UnmarshalText: %s", err)
		assert.Equal(t, start+5, no.Start, "NumberOption.UnmarshalText failed on start")
		assert.Equal(t, end+5, no.End, "NumberOption.UnmarshalText failed on end")

		err = no.UnmarshalText([]byte(fmt.Sprintf("%d-%d", start, end)))
		assert.NotNil(t, err, "NumberOption failed on UnmarshalText: %s", err)

		err = no.UnmarshalText([]byte(fmt.Sprintf("t%dy:%d", start, end)))
		assert.NotNil(t, err, "NumberOption failed on UnmarshalText: %s", err)
	})
}

var testNumberListOption = NumberListOption{
	Column: "table.column",
	Items: map[string][]interface{}{
		"1": {"a", "b", "c"},
		"2": {"d", "e", "f"},
		"3": {"g", "h", "i"},
	},
	Multiple: true,
	Joins: []JoinStatement{
		{
			Diff:    INNER_JOIN,
			Left:    JoinTable{"lkey", "lcol"},
			Right:   JoinTable{"rkey", "rcol"},
			GroupBy: "leky.lcol",
			OrderBy: "MAX(id)",
			Select:  "MAX(id) as last, MIN(id) as first",
		},
	},
}

func TestNumberListOption(t *testing.T) {
	t.Run("NumberListOption", func(t *testing.T) {
		no := testNumberListOption

		req := no.Request()
		assert.IsTypef(t, RequestList{}, req, "NumberListOption failed on zero-value RequestList: %T", req)

		rl, ok := req.(RequestList)
		assert.Truef(t, ok, "NumberOption failed on RequestList cast: %T", req)

		assert.Len(t, rl.list, 0, "NumberListOption failed on zero-value RequestList.list")
		assert.Len(t, rl.not, 9, "NumberListOption failed on zero-value RequestList.not")

		list := []string{"1", "3"}
		nm := no.Merge(NumberListOption{
			List: list,
		})

		assert.Lenf(t, nm.List, len(list), "NumberListOption.Merge failed on list: %d: %s || %d: %s", len(nm.List), nm.List, len(list), list)

		req = nm.Request()
		assert.IsTypef(t, RequestList{}, req, "NumberListOption failed on non-zero-value Request: %T", req)

		rl, ok = req.(RequestList)
		assert.Truef(t, ok, "NumberListOption failed on RequestList cast: %T", req)

		assert.Len(t, rl.list, 6, "NumberListOption failed on RequestList.list")
		assert.Len(t, rl.not, 3, "NumberListOption failed on RequestList.not")

		err := no.UnmarshalText([]byte(fmt.Sprintf("%d,%d,%d", 1, 2, 3)))
		assert.Nil(t, err, "NumberListOption failed on UnmarshalText: %s", err)
		assert.Len(t, no.List, 3, "NumberListOption.UnmarshalText failed on List len")
		assert.Equal(t, "1", no.List[0], "NumberListOption.UnmarshalText failed on List[0]")
		assert.Equal(t, "2", no.List[1], "NumberListOption.UnmarshalText failed on List[1]")
		assert.Equal(t, "3", no.List[2], "NumberListOption.UnmarshalText failed on List[2]")

		// err = no.UnmarshalText([]byte(fmt.Sprintf("%d-%d-%d", 1, 2, 3)))
		// assert.NotNil(t, err, "NumberListOption failed on UnmarshalText: %s", err)

		// err = no.UnmarshalText([]byte(fmt.Sprintf("t%dy,%dg,a %da", 1, 2, 3)))
		// assert.NotNil(t, err, "NumberListOption failed on UnmarshalText: %s", err)
	})
}
