package weather_gc_ca

import (
	"encoding/xml"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXML(t *testing.T) {
	t.Run("test hourly", func(t *testing.T) {
		b, err := ioutil.ReadFile("./_testdata/test-hourly_toronto.xml")
		if err != nil {
			t.Error(err)
		}
		d := &HourlyDataXML{}
		x := &ClimateDataXML{Data: d}
		err = xml.Unmarshal(b, x)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, "TORONTO LESTER B. PEARSON INT'L A", x.StationInfo.Name)
		assert.Equal(t, "ONTARIO", x.StationInfo.Province)
		assert.Equal(t, 43.68, x.StationInfo.Latitude)
		assert.Equal(t, -79.63, x.StationInfo.Longitude)
		assert.Equal(t, "6158733", x.StationInfo.ClimateID)
		assert.Equal(t, "71624", x.StationInfo.WMOID)
		assert.Equal(t, "YYZ", x.StationInfo.TCID)

		assert.Greater(t, len(*d), 0)
	})

	t.Run("test daily", func(t *testing.T) {
		b, err := ioutil.ReadFile("./_testdata/test-daily_toronto.xml")
		if err != nil {
			t.Error(err)
		}
		d := &DailyDataXML{}
		x := &ClimateDataXML{Data: d}
		err = xml.Unmarshal(b, x)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, "TORONTO LESTER B. PEARSON INT'L A", x.StationInfo.Name)
		assert.Equal(t, "ONTARIO", x.StationInfo.Province)
		assert.Equal(t, 43.68, x.StationInfo.Latitude)
		assert.Equal(t, -79.63, x.StationInfo.Longitude)
		assert.Equal(t, "6158733", x.StationInfo.ClimateID)
		assert.Equal(t, "71624", x.StationInfo.WMOID)
		assert.Equal(t, "YYZ", x.StationInfo.TCID)

		assert.Greater(t, len(*d), 0)
	})

	t.Run("test monthly", func(t *testing.T) {
		b, err := ioutil.ReadFile("./_testdata/test-monthly_toronto.xml")
		if err != nil {
			t.Error(err)
		}
		d := &MonthlyDataXML{}
		x := &ClimateDataXML{Data: d}
		err = xml.Unmarshal(b, x)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, "TORONTO LESTER B. PEARSON INT'L A", x.StationInfo.Name)
		assert.Equal(t, "ONTARIO", x.StationInfo.Province)
		assert.Equal(t, 43.68, x.StationInfo.Latitude)
		assert.Equal(t, -79.63, x.StationInfo.Longitude)
		assert.Equal(t, "6158733", x.StationInfo.ClimateID)
		assert.Equal(t, "71624", x.StationInfo.WMOID)
		assert.Equal(t, "YYZ", x.StationInfo.TCID)

		assert.Greater(t, len(*d), 0)
	})
}

func TestMethods(t *testing.T) {
	t.Run("Test Find", func(t *testing.T) {
		lat, lng := 50.4452, -104.6189
		s := StationInventory.Find(lat, lng, 25)
		assert.Equalf(t, len(s), 25, "expected 25 stations, got %d", len(s))
		for _, v := range s {
			// check that all stations are within 150km of the given lat/lng
			// 1.1degrees = ~122km
			adv := 1.1
			assert.LessOrEqualf(t, v.Latitude, lat+adv, "latitude for station %d is %.2f , expected <= %.2f", v.StationID, v.Latitude, lat)
			assert.GreaterOrEqualf(t, v.Latitude, lat-adv, "latitude for station %d is %.2f , expected >= %.2f", v.StationID, v.Latitude, lat)

			assert.LessOrEqualf(t, v.Longitude, lng+adv, "longitude for station %d is %.2f , expected <= %.2f", v.StationID, v.Longitude, lng)
			assert.GreaterOrEqualf(t, v.Longitude, lng-adv, "longitude for station %d is %.2f , expected >= %.2f", v.StationID, v.Longitude, lng)
		}
	})

	t.Run("Test Get", func(t *testing.T) {
		s, ok := StationInventory.Station(30247)
		if !ok {
			t.Error("could not find station")
		}
		assert.Equalf(t, s.StationID, 30247, "expected station id 30247, got %d", s.StationID)
		assert.Equalf(t, s.Name, "TORONTO CITY CENTRE", "expected station name TORONTO CITY CENTRE, got %s", s.Name)
	})
}
