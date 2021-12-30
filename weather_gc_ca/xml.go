package weather_gc_ca

import (
	"encoding/xml"
	"fmt"
	"sort"
	"time"
)

type ClimateDataXML struct {
	// b           []byte
	XMLName     xml.Name       `xml:"climatedata"`
	Lang        string         `xml:"lang" json:"lang"`
	StationInfo StationInfoXML `xml:"stationinformation" json:"stationinformation"`
	Legend      []FlagsXML     `xml:"legend>flag" json:"legend"`
	Data        StationDataXML `xml:"stationdata" json:"data"`
}

type StationDataXML interface {
	Timeframe() (Timeframe, Timeframe)
	Find(Timeframe) (IntervalBaseXML, bool)
	Sort()
	First() IntervalBaseXML
	Last() IntervalBaseXML
	Append(StationDataXML)
	csv() [][]string
}

type IntervalBaseXML interface {
	Timeframe() Timeframe
}

// TODO: almanac dataset

type MonthlyBaseXML struct {
	Time               time.Time `xml:"-" json:"time"`
	Month              int       `xml:"month,attr" json:"month"`
	Year               int       `xml:"year,attr" json:"year"`
	MeanMaxTemp        float64   `xml:"meanmaxtemp" json:"maxTemp"`
	MeanMinTemp        float64   `xml:"meanmintemp" json:"minTemp"`
	MeanTemp           float64   `xml:"meanmonthtemp" json:"meanTemp"`
	ExtremeMaxTemp     float64   `xml:"extrmaxtemp" json:"extremeMaxTemp"`
	ExtremeMinTemp     float64   `xml:"extrmintemp" json:"extremeMinTemp"`
	TotalRain          float64   `xml:"totrain" json:"rainfall"`
	TotalSnow          float64   `xml:"totsnow" json:"snowfall"`
	TotalPrecipitation float64   `xml:"totprecip" json:"totalPrecip"`
	SnowOnGround       float64   `xml:"grndsnowlastday" json:"snowDepth"`
	MaxGustDirection   float64   `xml:"dirmaxgust" json:"windDirection"`
	MaxGustSpeed       string    `xml:"speedmaxgust" json:"windGustSpeed"`
}

func (m MonthlyBaseXML) Timeframe() Timeframe {
	return Timeframe{
		Year:  m.Year,
		Month: m.Month,
		Time:  time.Date(m.Year, time.Month(m.Month), 1, 0, 0, 0, 0, time.UTC), // TODO: confirm TIMEZONE
	}
}

type MonthlyDataXML []MonthlyBaseXML

func (m *MonthlyDataXML) csv() [][]string {
	s := [][]string{}
	s = append(s, []string{
		"Year",
		"Month",
		"MeanMaxTemp",
		"MeanMinTemp",
		"MeanTemp",
		"ExtremeMaxTemp",
		"ExtremeMinTemp",
		"TotalRain",
		"TotalSnow",
		"TotalPrecipitation",
		"SnowOnGround",
		"MaxGustDirection",
		"MaxGustSpeed",
	})
	for _, a := range *m {
		s = append(s, []string{
			fmt.Sprintf("%d", a.Year),
			fmt.Sprintf("%d", a.Month),
			fmt.Sprintf("%.2f", a.MeanMaxTemp),
			fmt.Sprintf("%.2f", a.MeanMinTemp),
			fmt.Sprintf("%.2f", a.MeanTemp),
			fmt.Sprintf("%.2f", a.ExtremeMaxTemp),
			fmt.Sprintf("%.2f", a.ExtremeMinTemp),
			fmt.Sprintf("%.2f", a.TotalRain),
			fmt.Sprintf("%.2f", a.TotalSnow),
			fmt.Sprintf("%.2f", a.TotalPrecipitation),
			fmt.Sprintf("%.2f", a.SnowOnGround),
			fmt.Sprintf("%.2f", a.MaxGustDirection),
			a.MaxGustSpeed,
		})
	}

	return s
}

func (m *MonthlyDataXML) Timeframe() (start, end Timeframe) {
	m.Sort()
	dm := (*m)
	return dm[0].Timeframe(), dm[len(dm)-1].Timeframe()
}

func (m *MonthlyDataXML) Find(t Timeframe) (IntervalBaseXML, bool) {
	dm := (*m)
	for _, a := range dm {
		if a.Timeframe().Time.Equal(t.Time) {
			return a, true
		}
	}
	return nil, false
}

func (m *MonthlyDataXML) Sort() {
	md := (*m)
	sort.Slice(md, func(i, j int) bool {
		return md[i].Timeframe().Time.Before(md[j].Timeframe().Time)
	})
}

func (m *MonthlyDataXML) First() IntervalBaseXML {
	return (*m)[0]
}

func (m *MonthlyDataXML) Last() IntervalBaseXML {
	dm := (*m)
	return dm[len(dm)-1]
}

func (m *MonthlyDataXML) Append(data StationDataXML) {
	if v, ok := data.(*MonthlyDataXML); ok {
		dm := (*m)
		dv := (*v)
		dm = append(dm, dv...)
		*m = dm
	}
}

type DailyBaseXML struct {
	Time               time.Time `xml:"-" json:"time"`
	Day                int       `xml:"day,attr" json:"day"`
	Month              int       `xml:"month,attr" json:"month"`
	Year               int       `xml:"year,attr" json:"year"`
	MaxTemp            float64   `xml:"maxtemp" json:"maxTemp"`
	MinTemp            float64   `xml:"mintemp" json:"minTemp"`
	MeanTemp           float64   `xml:"meantemp" json:"meanTemp"`
	HeatDegDays        float64   `xml:"heatdegdays" json:"heatDegDays"`
	CoolDegDays        float64   `xml:"cooldegdays" json:"coolDegDays"`
	TotalRain          float64   `xml:"totalrain" json:"rainfall"`
	TotalSnow          float64   `xml:"totalsnow" json:"snowfall"`
	TotalPrecipitation float64   `xml:"totalprecipitation" json:"totalPrecip"`
	SnowOnGround       float64   `xml:"snowonground" json:"snowDepth"`
	MaxGustDirection   float64   `xml:"dirofmaxgust" json:"windDirection"`
	MaxGustSpeed       string    `xml:"speedofmaxgust" json:"windGustSpeed"`
}

func (d DailyBaseXML) Timeframe() Timeframe {
	return Timeframe{
		Year:  d.Year,
		Month: d.Month,
		Day:   d.Day,
		Time:  time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC), // TODO: confirm TIMEZONE
	}
}

type DailyDataXML []DailyBaseXML

func (d *DailyDataXML) csv() [][]string {
	s := [][]string{}
	s = append(s, []string{
		"Year",
		"Month",
		"Day",
		"MaxTemp",
		"MinTemp",
		"MeanTemp",
		"HeatDegDays",
		"CoolDegDays",
		"TotalRain",
		"TotalSnow",
		"TotalPrecipitation",
		"SnowOnGround",
		"MaxGustDirection",
		"MaxGustSpeed",
	})
	for _, a := range *d {
		s = append(s, []string{
			fmt.Sprintf("%d", a.Year),
			fmt.Sprintf("%d", a.Month),
			fmt.Sprintf("%d", a.Day),
			fmt.Sprintf("%.2f", a.MaxTemp),
			fmt.Sprintf("%.2f", a.MinTemp),
			fmt.Sprintf("%.2f", a.MeanTemp),
			fmt.Sprintf("%.2f", a.HeatDegDays),
			fmt.Sprintf("%.2f", a.CoolDegDays),
			fmt.Sprintf("%.2f", a.TotalRain),
			fmt.Sprintf("%.2f", a.TotalSnow),
			fmt.Sprintf("%.2f", a.TotalPrecipitation),
			fmt.Sprintf("%.2f", a.SnowOnGround),
			fmt.Sprintf("%.2f", a.MaxGustDirection),
			a.MaxGustSpeed,
		})
	}
	return s
}

func (d *DailyDataXML) Timeframe() (start, end Timeframe) {
	d.Sort()
	dd := (*d)
	return dd[0].Timeframe(), dd[len(dd)-1].Timeframe()
}

func (d *DailyDataXML) Find(t Timeframe) (IntervalBaseXML, bool) {
	dd := (*d)
	for _, a := range dd {
		if a.Timeframe().Time.Equal(t.Time) {
			return a, true
		}
	}
	return nil, false
}

func (d *DailyDataXML) Sort() {
	dd := (*d)
	sort.Slice(dd, func(i, j int) bool {
		return dd[i].Timeframe().Time.Before(dd[j].Timeframe().Time)
	})
}

func (d *DailyDataXML) First() IntervalBaseXML {
	return (*d)[0]
}

func (d *DailyDataXML) Last() IntervalBaseXML {
	dd := (*d)
	if len(dd) == 0 {
		return nil
	}
	return dd[len(dd)-1]
}

func (d *DailyDataXML) Append(data StationDataXML) {
	if v, ok := data.(*DailyDataXML); ok {
		dd := (*d)
		dv := (*v)
		dd = append(dd, dv...)
		*d = dd
	}
}

type HourlyBaseXML struct {
	Time             time.Time `xml:"-" json:"time"`
	Minute           int       `xml:"minute,attr" json:"minute"`
	Hour             int       `xml:"hour,attr" json:"hour"`
	Day              int       `xml:"day,attr" json:"day"`
	Month            int       `xml:"month,attr" json:"month"`
	Year             int       `xml:"year,attr" json:"year"`
	Temp             float64   `xml:"temp" json:"temp"`
	DewPointTemp     float64   `xml:"dptemp" json:"dewPointTemp"`
	RelativeHumidity float64   `xml:"relhum" json:"relativeHumidity"`
	WindDirection    float64   `xml:"winddir" json:"windDirection"`
	WindSpeed        string    `xml:"windspd" json:"windSpeed"`
	Visibility       float64   `xml:"visibility" json:"visibility"`
	StationPressure  float64   `xml:"stnpress" json:"stationPressure"`
	Humidex          float64   `xml:"humidex" json:"humidex"`
	Windchill        float64   `xml:"windchill" json:"windchill"`
	Weather          string    `xml:"weather" json:"weather"`
}

func (h HourlyBaseXML) Timeframe() Timeframe {
	return Timeframe{
		Year:  h.Year,
		Month: h.Month,
		Day:   h.Day,
		Time:  time.Date(h.Year, time.Month(h.Month), h.Day, h.Hour, h.Minute, 0, 0, time.UTC), // TODO: confirm TIMEZONE
	}
}

type HourlyDataXML []HourlyBaseXML

func (h *HourlyDataXML) csv() [][]string {
	s := [][]string{}
	s = append(s, []string{
		"Year",
		"Month",
		"Day",
		"Hour",
		"Minute",
		"Temp",
		"DewPointTemp",
		"RelativeHumidity",
		"WindDirection",
		"WindSpeed",
		"Visibility",
		"StationPressure",
		"Humidex",
		"Windchill",
		"Weather",
	})

	for _, a := range *h {
		s = append(s, []string{
			fmt.Sprintf("%d", a.Year),
			fmt.Sprintf("%d", a.Month),
			fmt.Sprintf("%d", a.Day),
			fmt.Sprintf("%d", a.Hour),
			fmt.Sprintf("%d", a.Minute),
			fmt.Sprintf("%.2f", a.Temp),
			fmt.Sprintf("%.2f", a.DewPointTemp),
			fmt.Sprintf("%.2f", a.RelativeHumidity),
			fmt.Sprintf("%.2f", a.WindDirection),
			a.WindSpeed,
			fmt.Sprintf("%.2f", a.Visibility),
			fmt.Sprintf("%.2f", a.StationPressure),
			fmt.Sprintf("%.2f", a.Humidex),
			fmt.Sprintf("%.2f", a.Windchill),
			a.Weather,
		})
	}

	return s
}

func (h *HourlyDataXML) Timeframe() (start, end Timeframe) {
	h.Sort()
	hd := (*h)
	return hd[0].Timeframe(), hd[len(hd)-1].Timeframe()
}

func (h *HourlyDataXML) Find(t Timeframe) (IntervalBaseXML, bool) {
	hd := (*h)
	for _, a := range hd {
		if a.Timeframe().Time.Equal(t.Time) {
			return a, true
		}
	}
	return nil, false
}

func (h *HourlyDataXML) Sort() {
	hd := (*h)
	sort.Slice(hd, func(i, j int) bool {
		return hd[i].Timeframe().Time.Before(hd[j].Timeframe().Time)
	})
}

func (h *HourlyDataXML) First() IntervalBaseXML {
	return (*h)[0]
}

func (h *HourlyDataXML) Last() IntervalBaseXML {
	hd := (*h)
	return hd[len(hd)-1]
}

func (h *HourlyDataXML) Append(data StationDataXML) {
	if v, ok := data.(*HourlyDataXML); ok {
		hd := (*h)
		dv := (*v)
		hd = append(hd, dv...)
		*h = hd
	}
}

type StationInfoXML struct {
	Name      string  `xml:"name" json:"name"`
	Province  string  `xml:"province" json:"province"`
	Latitude  float64 `xml:"latitude" json:"latitude"`
	Longitude float64 `xml:"longitude" json:"longitude"`
	Elevation float64 `xml:"elevation" json:"elevation"`
	ClimateID string  `xml:"climate_identifier" json:"climateID"`
	WMOID     string  `xml:"wmo_identifier" json:"wmoID"`
	TCID      string  `xml:"tc_identifier" json:"tcID"`
}

type FlagsXML struct {
	Symbol      string `xml:"symbol" json:"symbol"`
	Description string `xml:"description" json:"description"`
}
