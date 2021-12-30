package weather_gc_ca

import (
	"encoding/json"
	"errors"
	"time"
)

type RawStations []StationMetadata

type StationMetadata struct {
	previousDistance float64

	XML              ClimateDataXML `xml:"-" json:"-"`
	Name             string         `json:"Name"`
	Province         string         `json:"Province"`
	ClimateID        string         `json:"Climate ID"`
	StationID        int            `json:"Station ID"`
	WMOID            string         `json:"WMO ID"`
	TCID             string         `json:"TC ID"`
	Latitude         float64        `json:"Latitude (Decimal Degrees)"`
	Longitude        float64        `json:"Longitude (Decimal Degrees)"`
	Latitude_int     int32          `json:"Latitude"`
	Longitude_int    int32          `json:"Longitude"`
	Elevation        float64        `json:"Elevation (m)"`
	FirstYear        int            `json:"First Year"`
	LastYear         int            `json:"Last Year"`
	HourlyFirstYear  int            `json:"HLY First Year"`
	HourlyLastYear   int            `json:"HLY Last Year"`
	DailyFirstYear   int            `json:"DLY First Year"`
	DailyLastYear    int            `json:"DLY Last Year"`
	MonthlyFirstYear int            `json:"MLY First Year"`
	MonthlyLastYear  int            `json:"MLY Last Year"`
}

func (s StationMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"name":             s.Name,
		"stationID":        s.StationID,
		"province":         s.Province,
		"latitude":         s.Latitude,
		"longitude":        s.Longitude,
		"elevation":        s.Elevation,
		"firstYear":        s.FirstYear,
		"lastYear":         s.LastYear,
		"hourlyFirstYear":  s.HourlyFirstYear,
		"hourlyLastYear":   s.HourlyLastYear,
		"dailyFirstYear":   s.DailyFirstYear,
		"dailyLastYear":    s.DailyLastYear,
		"monthlyFirstYear": s.MonthlyFirstYear,
		"monthlyLastYear":  s.MonthlyLastYear,
	})
}

type Interval int

func (i Interval) String() string {
	switch i {
	case 1:
		return "Hourly"
	case 2:
		return "Daily"
	case 3:
		return "Monthly"
	case 4:
		return "Almanac"
	}
	return "Unknown"
}

const (
	Hourly  Interval = 1
	Daily   Interval = 2
	Monthly Interval = 3
	Almanac Interval = 4
)

type Timeframe struct {
	Day   int       `json:"day"`
	Month int       `json:"month"`
	Year  int       `json:"year"`
	Time  time.Time `json:"time"`
}

func (t Timeframe) String() string {
	if t.Time.IsZero() {
		t.Time = time.Date(t.Year, time.Month(t.Month), t.Day, 0, 0, 0, 0, time.UTC)
		return t.Time.Format("01/02/06")
	}
	return t.Time.Format("01/02/06 15:04:05") // Mon Jan 2 15:04:05 -0700 MST 2006
}

type DownloadStatus struct {
	Progress chan DownloadProgress
	Done     chan bool
}

type DownloadProgress struct {
	Timeframe Timeframe
	Total     int
	Count     int
	Time      int64
	Error     error
}

var (
	ErrContextCancelled = errors.New("context cancelled")
	ErrRequestFailed    = errors.New("request failed")
)
