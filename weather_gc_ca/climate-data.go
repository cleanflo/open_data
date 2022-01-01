package weather_gc_ca

import (
	"context"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

func (r RawStations) String() string {
	a := "Distance\tID\t\tHourly\t\tDaily\t\tMonthly\t\tName\n"
	for _, s := range r {
		a += fmt.Sprintf("%.2f\tkm\t%d\t\t%d\t%d\t%d\t%d\t%d\t%d\t%s\n",
			s.previousDistance, s.StationID,
			s.HourlyFirstYear, s.HourlyLastYear,
			s.DailyFirstYear, s.DailyLastYear,
			s.MonthlyFirstYear, s.MonthlyLastYear,
			s.Name)
	}

	return a
}

func (r RawStations) Map(f func(a *StationMetadata)) {
	for _, a := range r {
		f(&a)
	}
}

type distanceMap struct {
	list     map[float64]StationMetadata
	farthest float64
}

func (m *distanceMap) add(d float64, s StationMetadata, max int) {
	// fill map with first n[max] stations on list
	if len(m.list) < max {
		m.list[d] = s
		if d > m.farthest {
			m.farthest = d
		}
		return
	}

	// if this station is closer than the farthest station in the map
	if d < m.farthest {
		// remove farthest station from map
		delete(m.list, m.farthest)
		// add this station to map
		m.list[d] = s
		m.farthest = 0.0
		// find the new farthest station in the map
		for distance := range m.list {
			if distance > m.farthest {
				m.farthest = distance
			}
		}
	}
}

func (r RawStations) Find(lat, lng float64, max int) (s RawStations) {
	// get a list of stations sorted by distance
	m := distanceMap{list: make(map[float64]StationMetadata, max)}

	// TODO: figure out how to start excluding stations early
	for _, a := range r {
		d := a.Distance(lat, lng)
		m.add(d, a, max)
		// // fill map with first n[max] stations on list
		// if len(m) < max {
		// 	m[d] = a
		// 	if d > farthest {
		// 		farthest = d
		// 	}
		// 	continue
		// }

		// // if this station is closer than the farthest station in the map
		// if d < farthest {
		// 	// remove farthest station from map
		// 	delete(m, farthest)
		// 	// add this station to map
		// 	m[d] = a
		// 	farthest = 0.0
		// 	// find the new farthest station in the map
		// 	for distance := range m {
		// 		if distance > farthest {
		// 			farthest = distance
		// 		}
		// 	}
		// }
	}

	s = nil
	// convert map to slice
	for _, v := range m.list {
		s = append(s, v)
	}

	return s
}

func (r RawStations) FindWithInterval(lat, lng float64, max int, interval Interval) (s RawStations) {
	// get a list of stations sorted by distance
	m := &distanceMap{list: make(map[float64]StationMetadata, max)}
	for _, a := range r {
		switch interval {
		case Hourly:
			if a.HourlyFirstYear == 0 || a.HourlyLastYear == 0 {
				continue
			}
		case Daily:
			if a.DailyFirstYear == 0 || a.DailyLastYear == 0 {
				continue
			}
		case Monthly:
			if a.MonthlyFirstYear == 0 || a.MonthlyLastYear == 0 {
				continue
			}
		case Almanac:
			if a.FirstYear == 0 || a.LastYear == 0 {
				continue
			}
		}
		d := a.Distance(lat, lng)
		m.add(d, a, max)
		// fill map with first n[max] stations on list
		// if len(m) < max {
		// 	m[d] = a
		// 	if d > farthest {
		// 		farthest = d
		// 	}
		// 	continue
		// }

		// // if this station is closer than the farthest station in the map
		// if d < farthest {
		// 	// remove farthest station from map
		// 	delete(m, farthest)
		// 	// add this station to map
		// 	m[d] = a
		// 	farthest = 0.0
		// 	// find the new farthest station in the map
		// 	for distance := range m {
		// 		if distance > farthest {
		// 			farthest = distance
		// 		}
		// 	}
		// }
	}

	s = nil
	// convert map to slice
	for _, v := range m.list {
		s = append(s, v)
	}

	return s
}

func (r RawStations) NameContains(query string, max int) (s RawStations) {
	for _, a := range r {
		if strings.Contains(strings.ToLower(a.Name), strings.ToLower(query)) {
			s = append(s, a)
			if len(s) >= max {
				return s
			}
		}
	}
	return s
}

func (r RawStations) NameStartsWith(query string, max int) (s RawStations) {
	for _, a := range r {
		if strings.HasPrefix(strings.ToLower(a.Name), strings.ToLower(query)) {
			s = append(s, a)
			if len(s) >= max {
				return s
			}
		}
	}
	return s
}

func (r RawStations) Station(id int) (s StationMetadata, ok bool) {
	for _, a := range r {
		if a.StationID == id {
			return a, true
		}
	}
	return s, false
}

func (r RawStations) csv() (a [][]string) {
	a = append(a, []string{
		"Station ID",
		"Name",
		"Province",
		"Climate ID",
		"WMO ID",
		"TC ID",
		"Latitude",
		"Longitude",
		"Elevation",
		"First Year",
		"Last Year",
		"Hourly First Year",
		"Hourly Last Year",
		"Daily First Year",
		"Daily Last Year",
		"Monthly First Year",
		"Monthly Last Year",
	})
	for _, s := range r {
		a = append(a, []string{
			fmt.Sprintf("%d", s.StationID),
			s.Name,
			s.Province,
			s.ClimateID,
			s.WMOID,
			s.TCID,
			fmt.Sprintf("%.2f", s.Latitude),
			fmt.Sprintf("%.2f", s.Longitude),
			fmt.Sprintf("%.2f", s.Elevation),
			fmt.Sprintf("%d", s.FirstYear),
			fmt.Sprintf("%d", s.LastYear),
			fmt.Sprintf("%d", s.HourlyFirstYear),
			fmt.Sprintf("%d", s.HourlyLastYear),
			fmt.Sprintf("%d", s.DailyFirstYear),
			fmt.Sprintf("%d", s.DailyLastYear),
			fmt.Sprintf("%d", s.MonthlyFirstYear),
			fmt.Sprintf("%d", s.MonthlyLastYear),
		})
	}
	return a
}

func (r RawStations) CSV(w io.Writer) error {
	return csv.NewWriter(w).WriteAll(r.csv())
}

type SortBy int

const (
	SortByDistance SortBy = iota
	SortByName
	SortByHourly
	SortByDaily
	SortByMonthly
)

func SortByString(a string) SortBy {
	switch strings.ToLower(a) {
	case "distance":
		return SortByDistance
	case "name":
		return SortByName
	case "hourly":
		return SortByHourly
	case "daily":
		return SortByDaily
	case "monthly":
		return SortByMonthly
	default:
		return SortByName
	}
}

func (r RawStations) Sort(by SortBy) {
	switch by {
	case SortByDistance:
		sort.Slice(r, func(i, j int) bool {
			return r[i].previousDistance < r[j].previousDistance
		})
	case SortByName:
		sort.Slice(r, func(i, j int) bool {
			return r[i].Name < r[j].Name
		})
	case SortByHourly:
		sort.Slice(r, func(i, j int) bool {
			return (r[i].HourlyLastYear - r[i].HourlyFirstYear) < (r[j].HourlyLastYear - r[j].HourlyFirstYear)
		})
	case SortByDaily:
		sort.Slice(r, func(i, j int) bool {
			return (r[i].DailyLastYear - r[i].DailyFirstYear) < (r[j].DailyLastYear - r[j].DailyFirstYear)
		})
	case SortByMonthly:
		sort.Slice(r, func(i, j int) bool {
			return (r[i].MonthlyLastYear - r[i].MonthlyFirstYear) < (r[j].MonthlyLastYear - r[j].MonthlyFirstYear)
		})
	}
}

func (r RawStations) Distances() map[float64]StationMetadata {
	m := make(map[float64]StationMetadata)
	for _, a := range r {
		m[a.previousDistance] = a
	}
	return m
}

func (r *StationMetadata) String() string {
	return fmt.Sprintf(
		`Station ID:	%d
Name:		%s
Province:	%s
Latitude:	%.2f
Longitude:	%.2f
Elevation:	%.2f m
Hourly:		%d - %d
Daily:		%d - %d
Monthly:	%d - %d
`, r.StationID, r.Name, r.Province, r.Latitude, r.Longitude, r.Elevation,
		r.HourlyFirstYear, r.HourlyLastYear, r.DailyFirstYear, r.DailyLastYear, r.MonthlyFirstYear, r.MonthlyLastYear)

}

func (r *StationMetadata) CSV(w io.Writer) error {
	return csv.NewWriter(w).WriteAll(r.XML.Data.csv())
}

func (r *StationMetadata) Timeframe(interval Interval) (start int, end int) {
	switch interval {
	case Hourly:
		return r.HourlyFirstYear, r.HourlyLastYear
	case Daily:
		return r.DailyFirstYear, r.DailyLastYear
	case Monthly:
		return r.MonthlyFirstYear, r.MonthlyLastYear
	}
	return 0, 0
}

func (r *StationMetadata) Distance(lat, lng float64) float64 {
	rad := 6371.0
	dlat := r.radians(r.Latitude - lat)
	dlng := r.radians(r.Longitude - lng)

	n := (math.Pow(math.Sin(dlat/2), 2) + math.Cos(r.radians(lat))*math.Cos(r.radians(r.Latitude))*math.Pow(math.Sin(dlng/2), 2))
	angle := 2 * math.Asin(math.Sqrt(n))

	r.previousDistance = angle * rad

	return r.previousDistance
}

func (r StationMetadata) radians(deg float64) float64 {
	return deg * math.Pi / 180
}

// https://climate.weather.gc.ca/climate_data/bulk_data_e.html
//				?format=xml&stationID=5097&Year=${year}&Month=${month}&Day=1&timeframe=2&submit= Download+Data
func (r *StationMetadata) RetreiveData(year, month, day int, interval Interval) error {
	q := url.Values{}
	q.Add("format", "xml")
	q.Add("stationID", fmt.Sprintf("%d", r.StationID))
	q.Add("Year", fmt.Sprintf("%d", year))
	q.Add("Month", fmt.Sprintf("%d", month))
	q.Add("Day", fmt.Sprintf("%d", day))
	q.Add("timeframe", fmt.Sprintf("%d", interval))
	q.Add("submit", " Download+Data")
	u := url.URL{
		Scheme:   "https",
		Host:     "climate.weather.gc.ca",
		Path:     "/climate_data/bulk_data_e.html",
		RawQuery: q.Encode(),
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to retreive dataset: %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to retreive dataset: %s", err)
	}

	x := ClimateDataXML{}
	switch interval {
	case Hourly:
		x.Data = &HourlyDataXML{}
	case Daily:
		x.Data = &DailyDataXML{}
	case Monthly:
		x.Data = &MonthlyDataXML{}
	}
	err = xml.NewDecoder(resp.Body).Decode(&x)
	if err != nil {
		return fmt.Errorf("failed to retreive dataset: %s", err)
	}

	r.XML.Data.Append(x.Data)

	return nil
}

func (r *StationMetadata) RetreiveHourlyData(ctx context.Context) DownloadStatus {
	start, end := Timeframe{
		Day:   1,
		Month: 1,
		Year:  r.HourlyFirstYear,
	}, Timeframe{
		Day:   31,
		Month: 12,
		Year:  r.HourlyLastYear,
	}

	return r.retreiveBetween(ctx, start, end, Hourly)
}

func (r *StationMetadata) RetreiveDailyData(ctx context.Context) DownloadStatus {
	start, end := Timeframe{
		Day:   1,
		Month: 1,
		Year:  r.DailyFirstYear,
	}, Timeframe{
		Day:   31,
		Month: 12,
		Year:  r.DailyLastYear,
	}
	return r.retreiveBetween(ctx, start, end, Daily)
}

func (r *StationMetadata) RetreiveMonthlyData(ctx context.Context) DownloadStatus {
	start, end := Timeframe{
		Day:   1,
		Month: 1,
		Year:  r.MonthlyFirstYear,
	}, Timeframe{
		Day:   31,
		Month: 12,
		Year:  r.MonthlyLastYear,
	}
	return r.retreiveBetween(ctx, start, end, Monthly)
}

func (r *StationMetadata) RetreiveInterval(ctx context.Context, interval Interval) DownloadStatus {
	start, end := r.Timeframe(interval)
	return r.retreiveBetween(ctx,
		Timeframe{Year: start, Month: 1, Day: 1},
		Timeframe{Year: end, Month: 12, Day: 31},
		interval,
	)
}

func (r *StationMetadata) RetreiveTimeframe(ctx context.Context, start, end Timeframe, interval Interval) DownloadStatus {
	return r.retreiveBetween(ctx, start, end, interval)
}

func (r *StationMetadata) retreiveBetween(ctx context.Context, start, end Timeframe, interval Interval) DownloadStatus {
	yr, eyr, mon, emon := start.Year, end.Year, 1, 12

	switch interval {
	case Hourly:
		// yr = r.HourlyFirstYear
		// eyr = r.HourlyLastYear
		r.XML.Data = &HourlyDataXML{}
	case Daily:
		// yr = r.DailyFirstYear
		// eyr = r.DailyLastYear
		emon = 1
		r.XML.Data = &DailyDataXML{}
	case Monthly:
		// yr = r.MonthlyFirstYear
		// eyr = r.MonthlyFirstYear // first year because all data comes in 1 file
		emon = 1
		r.XML.Data = &MonthlyDataXML{}
	default:
		return DownloadStatus{}
	}

	total := (eyr - yr) * (emon - mon + 1)

	d := DownloadStatus{
		Progress: make(chan DownloadProgress),
		Done:     make(chan bool),
	}

	go func() {
		// iterate each year, month
		count := 1
		for ; yr < eyr; yr++ {
		M:
			for ; mon <= emon; mon++ {
				// try to find the data in the existing data
				// if hourly, we assume the entire month exists if the first entry exists
				// if daily/monthly, we assume the entire year exists if the first entry exists
				if _, ok := r.XML.Data.Find(Timeframe{
					Year:  yr,
					Month: mon,
					Day:   1,
					Time:  time.Date(yr, time.Month(mon), 1, 0, 0, 0, 0, time.UTC),
				}); ok {
					continue M
				}

				// download the data
				select {
				case <-ctx.Done():
					d.Done <- false
					return
				default:
					progress := DownloadProgress{
						Total: total,
						Count: count,
					}

					err := r.RetreiveData(yr, mon, 1, interval)
					if err != nil {
						fmt.Println(err)
						progress.Error = ErrRequestFailed
						d.Progress <- progress
						continue M
					}

					progress.Timeframe = r.XML.Data.Last().Timeframe()
					progress.Time = time.Now().Unix()
					d.Progress <- progress

					count++
				}
			}
			mon = 1
		}

		d.Done <- true
	}()

	return d
}
