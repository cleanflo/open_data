package db_layer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRequestRange(t *testing.T) {
	t.Run("RequestRange:TimeOption", func(t *testing.T) {
		to := testTimeOption

		start, end := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC), time.Now()
		tm := to.Merge(TimeOption{Start: start, End: end})

		req := tm.Request()
		assert.IsTypef(t, RequestRange{}, req, "TimeOption failed on Request: %T", req)
		rr, ok := req.(RequestRange)
		assert.Truef(t, ok, "TimeOption failed on RequestRange cast: %T", req)

		q, err := req.Query("testTimeColumn")
		assert.Nilf(t, err, "TimeOption failed on Query: %s", err)

		assert.Equal(t, "testTimeColumn BETWEEN ? AND ?", q.query, "TimeOption failed on Query")
		assert.Equal(t, []interface{}{rr.Start, rr.End}, q.value, "TimeOption failed on Query")

		rr.End = nil
		q, err = rr.Query("testTimeColumn")
		assert.Nilf(t, err, "TimeOption failed on Query: %s", err)

		assert.Equal(t, "testTimeColumn >= ?", q.query, "TimeOption failed on Query")
		assert.Equal(t, rr.Start, q.value, "TimeOption failed on Query")

		rr.End = rr.Start
		rr.Start = nil
		q, err = rr.Query("testTimeColumn")
		assert.Nilf(t, err, "TimeOption failed on Query: %s", err)

		assert.Equal(t, "testTimeColumn <= ?", q.query, "TimeOption failed on Query")
		assert.Equal(t, rr.End, q.value, "TimeOption failed on Query")

		rr.End = nil
		q, err = rr.Query("testTimeColumn")
		assert.NotNil(t, err, "TimeOption failed on Query: %s", err)
	})

	t.Run("RequestRange:NumberOption", func(t *testing.T) {
		no := testNumberOption

		start, end := 50, 500
		nm := no.Merge(NumberOption{Start: start, End: end})

		req := nm.Request()
		assert.IsTypef(t, RequestRange{}, req, "NumberOption failed on Request: %T", req)
		rr, ok := req.(RequestRange)
		assert.Truef(t, ok, "NumberOption failed on RequestRange cast: %T", req)

		q, err := req.Query("testNumColumn")
		assert.Nilf(t, err, "NumberOption failed on Query: %s", err)

		assert.Equal(t, "testNumColumn BETWEEN ? AND ?", q.query, "NumberOption failed on Query")
		assert.Equal(t, []interface{}{rr.Start, rr.End}, q.value, "NumberOption failed on Query")

		rr.End = nil
		q, err = rr.Query("testNumColumn")
		assert.Nilf(t, err, "NumberOption failed on Query: %s", err)

		assert.Equal(t, "testNumColumn >= ?", q.query, "NumberOption failed on Query")
		assert.Equal(t, rr.Start, q.value, "NumberOption failed on Query")

		rr.End = rr.Start
		rr.Start = nil
		q, err = rr.Query("testNumColumn")
		assert.Nilf(t, err, "NumberOption failed on Query: %s", err)

		assert.Equal(t, "testNumColumn <= ?", q.query, "NumberOption failed on Query")
		assert.Equal(t, rr.End, q.value, "NumberOption failed on Query")

		rr.End = nil
		q, err = rr.Query("testNumColumn")
		assert.NotNil(t, err, "NumberOption failed on Query: %s", err)
	})
}

func TestRequestList(t *testing.T) {
	t.Run("RequestList", func(t *testing.T) {
		no := testNumberListOption

		list := []string{"1"}
		nm := no.Merge(NumberListOption{
			List: list,
		})

		req := nm.Request()
		assert.IsTypef(t, RequestList{}, req, "NumberListOption failed on non-zero-value Request: %T", req)

		rl, ok := req.(RequestList)
		assert.Truef(t, ok, "NumberListOption failed on RequestList cast: %T", req)

		q, err := rl.Query("testNumColumn")
		assert.Nilf(t, err, "NumberListOption failed on Query: %s", err)

		assert.Equal(t, "testNumColumn IN ?", q.query, "NumberListOption failed on Query")
		assert.Equal(t, rl.list, q.value, "NumberListOption failed on Query")

		nm.List = append(nm.List, "2")
		req = nm.Request()
		rl = req.(RequestList)

		q, err = req.Query("testNumColumn")
		assert.Nilf(t, err, "NumberListOption failed on Query: %s", err)

		assert.Equal(t, "testNumColumn NOT IN ?", q.query, "NumberListOption failed on Query")
		assert.Equal(t, rl.not, q.value, "NumberListOption failed on Query")
	})
}
