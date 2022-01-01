package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	climatedata "github.com/cleanflo/open_data/weather_gc_ca"
)

func init() {

}

func main() {
	app := &cli.App{
		Name:  "climate data (climate.weather.gc.ca) CLI",
		Usage: "a tool for downloading climate data from the government of canada",
		Description: `
1. Search for a station:
	climate search distance --lat 45 --lon 123 --max 25
OR
	climate search name --contains calgary --max 5

2. Download data:
	climate download --stn 1234 --interval daily --start 1970 --end 2021
`,
		Commands: []*cli.Command{
			{
				Name: "search",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "max-count",
						Aliases: []string{"m", "max"},
						Value:   20,
						Usage:   "max number of stations to return",
					},
					&cli.StringFlag{
						Name:    "sort",
						Aliases: []string{"s"},
						Usage:   "sort by field: distance, name, hourly, daily, monthly",
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "search for stations by name",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "contains",
								Usage: "search for stations containing this string",
							},
							&cli.StringFlag{
								Name:  "starts-with",
								Usage: "search for stations starting with this string",
							},
						},
						Action: SearchByName,
					},
					{
						Name:    "distance",
						Aliases: []string{"d"},
						Usage:   "output a list of stations from a given coordinate pair",
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:     "latitude",
								Aliases:  []string{"lat"},
								Usage:    "Latitude of the coordinate pair",
								Required: true,
							},
							&cli.Float64Flag{
								Name:     "longitude",
								Aliases:  []string{"lon"},
								Usage:    "Longitude of the coordinate pair",
								Required: true,
							},
							&cli.IntFlag{
								Name:    "interval",
								Aliases: []string{"int", "i"},
								Usage:   "The desired interval of data: hourly(0), daily(1), monthly(2)",
							},
						},
						Action: SearchByCoor,
					},
				},
			},
			{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "retreive info for a station by ID",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "station id",
						Aliases:  []string{"s", "stn", "id"},
						Usage:    "Station ID to get info for",
						Required: true,
					},
				},
				Action: StationInfo,
			},
			{
				Name:  "download",
				Usage: "download data for a station, if no start or end is supplied it will download the entire time range",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "station id",
						Aliases:  []string{"s", "stn", "id"},
						Usage:    "Station ID to get info for",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o", "f", "file"},
						Usage:   "File to write output of successful download `FILE`",
					},
					&cli.StringFlag{
						Name:    "interval",
						Aliases: []string{"i", "int"},
						Value:   "daily",
						Usage:   "interval to download data for: hourly, daily, monthly",
					},
					&cli.IntFlag{
						Name:  "start",
						Usage: "starting year to download data for the station",
					},
					&cli.IntFlag{
						Name:  "end",
						Usage: "ending year to download data for the station",
					},
				},
				Action: DownloadData,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func SearchByName(c *cli.Context) error {
	contains := c.String("contains")
	startsWith := c.String("starts-with")
	max := c.Int("max-count")

	if contains != "" && startsWith != "" {
		return fmt.Errorf("cannot specify both --contains and --starts-with")
	}

	if contains == "" && startsWith == "" {
		return fmt.Errorf("must specify either --contains or --starts-with")
	}

	var stations climatedata.RawStations
	if contains != "" {
		stations = climatedata.StationInventory.NameContains(contains, max)
	} else if startsWith != "" {
		stations = climatedata.StationInventory.NameStartsWith(startsWith, max)
	}
	stations.Sort(climatedata.SortByString(c.String("sort")))
	fmt.Println(stations)

	return nil
}

func SearchByCoor(c *cli.Context) error {
	lat := c.Float64("latitude")
	lng := c.Float64("longitude")
	max := c.Int("max-count")
	interval := climatedata.Interval(c.Int("interval"))
	stations := climatedata.StationInventory.FindWithInterval(lat, lng, max, interval)

	sort := c.String("sort")
	if sort == "" {
		sort = "distance"
	}

	stations.Sort(climatedata.SortByString(sort))
	fmt.Println(stations)

	return nil
}

func StationInfo(c *cli.Context) error {
	stn := c.Int("station")
	station, ok := climatedata.StationInventory.Station(stn)
	if !ok {
		return fmt.Errorf("station %d not found", stn)
	}

	fmt.Println(station)
	return nil
}
