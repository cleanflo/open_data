package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	climatedata "github.com/cleanflo/open_data/weather_gc_ca"
	"github.com/urfave/cli/v2"
)

func DownloadData(c *cli.Context) error {
	var interval climatedata.Interval
	switch c.String("interval") {
	case "hourly":
		interval = climatedata.Hourly
	case "daily":
		interval = climatedata.Daily
	case "monthly":
		interval = climatedata.Monthly
	default:
		return fmt.Errorf("invalid interval: %s", c.String("interval"))
	}

	stn := c.Int("station")
	start := climatedata.Timeframe{
		Year:  c.Int("start"),
		Month: 1,
		Day:   1,
	}
	end := climatedata.Timeframe{
		Year:  c.Int("end"),
		Month: 12,
		Day:   31,
	}

	s, ok := climatedata.StationInventory.Station(stn)
	if !ok {
		return fmt.Errorf("station %d not found", stn)
	}

	startYear, endYear := s.Timeframe(interval)
	if start.Year == 0 {
		start.Year = startYear
	} else {
		if start.Year < startYear {
			return fmt.Errorf("provided start year %d is before station start year %d", start, startYear)
		}
	}
	if end.Year == 0 {
		end.Year = endYear
	} else {
		if end.Year > endYear {
			return fmt.Errorf("provided end year %d is after station end year %d", end, endYear)
		}
	}

	p := c.Path("output")
	if p == "" {
		p = fmt.Sprintf("./%s_%d_%s_%d-%d.csv", s.Name, s.StationID, interval, start.Year, end.Year)
	}

	outputFile, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	defer outputFile.Close()

	fmt.Printf("Downloading %s data for station %d from %s to %s\n", interval, stn, start, end)

	childCtx, cancel := context.WithCancel(c.Context)
	progress := s.RetreiveTimeframe(childCtx, start, end, interval)

	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(ch, os.Interrupt)
	for {
		select {
		case <-ch:
			fmt.Println("Received SIGINT, stopping...")
			cancel()
			os.Exit(0)
		case <-c.Context.Done():
			fmt.Println("context cancelled")
			return nil
		case p := <-progress.Progress:
			if p.Error != nil {
				fmt.Printf("Error: %s\n", p.Error)
				continue
			}
			fmt.Printf("Finished downloading data (%d / %d) upto: %s %d\n", p.Count, p.Total, time.Month(p.Timeframe.Month), p.Timeframe.Year)
		case t := <-progress.Done:
			if t {
				fmt.Println("Download complete")
			} else {
				fmt.Println("Download did not complete!")
			}

			err = s.CSV(outputFile)
			if err != nil {
				return fmt.Errorf("failed to write CSV: %w", err)
			}
			fmt.Println("CSV written to", p)
			return nil
		default:
			time.Sleep(time.Millisecond * 200)
		}
	}
}
