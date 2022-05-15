# Water Well Data

This package provides a simple API/CLI for searching the Environment Canada "Weather Station Inventory". The published data was collected in CSV format and converted to JSON for embedding in the package. The station inventory can be searched using the CLI or by using the provided `net/http.HandlerFunc` to serve the API. Once a "StationID" is found, the data can be retrieved using the provided options in the CLI or endpoint in the API, currently the data is concatonated into a single CSV file:

### Data Sources
[Technical Documentation](https://climate.weather.gc.ca/doc/Technical_Documentation.pdf)

- [Station Inventory (Google Drive)](https://drive.google.com/file/d/1HDRnj41YBWpMioLPwAFiLlK4SK8NV72C/view?usp=sharing)
- [Licence](https://climate.weather.gc.ca/prods_servs/attachment1_e.html)


## Usage
The first means of using the API is via the `Handler` method, which is a `net/http.HandlerFunc` that handles the request, performs the query, and returns the data. The `Handler` method can be used as below:

```go
dataRouter := r.PathPrefix("/data/").Subrouter()

weatherRouter := dataRouter.PathPrefix("/weather/").Subrouter()
weatherRouter.HandleFunc("/station/search/", climatedata.SearchHandler).Methods("GET")
weatherRouter.HandleFunc("/station/download/", climatedata.DownloadHandler).Methods("GET")

```

You can also use the CLI to search the Station Inventory and download the data. The CLI can be used as below:

```bash
~$ ./climate-data search distance --lat 45.5 --lon -75.5
Distance        ID              Hourly          Daily           Monthly         Name
3.90    km      4261            0       0       1973    1980    1973    1980    CUMBERLAND
5.18    km      53001           2014    2021    2018    2021    0       0       OTTAWA GATINEAU A
5.89    km      8375            1981    2012    0       0       0       0       GATINEAU A
6.79    km      5574            0       0       1962    2021    1962    2018    ANGERS
7.07    km      5610            0       0       1963    1979    1963    1979    MASSON
...........

~$ ./climate-data search name --starts-with calgary
Distance        ID              Hourly          Daily           Monthly         Name
0.00    km      2202            0       0       1965    1965    1965    1965    CALGARY BEARSPAW
0.00    km      2203            0       0       1961    1966    1961    1966    CALGARY BELLEVIEW
0.00    km      2168            0       0       1966    1987    1966    1987    CALGARY ELBOW VIEW
0.00    km      2204            0       0       1956    1979    1956    1979    CALGARY GLENMORE DAM
0.00    km      2205            1953    2012    1881    2012    1881    2012    CALGARY INT\'L A
...........

~$ ./climate-data download --stn 2203 --start 1965 --interval daily
Downloading Daily data for station 2203 from 01/01/65 to 12/31/66
Finished downloading data (1 / 1) upto: December 1965
Download complete
CSV written to ./CALGARY BELLEVIEW_2203_Daily_1965-1966.csv
```

## Description
This package attempts to abstract the query logic and provide a simple http endpoint for searching the Environment Canada Station Inventory and querying the Environment Canada API to download the data. The requests are made using GET parameters to enable response caching on a variety of host providers, and the response is returned as a JSON. 


### Search Options
 - Global Flags:
   - max-count: The maximum number of results to return.
   - sort: sort by field: distance, name, hourly, daily, monthly
- Search
  - name: search by station name
    - contains: search for partial station name
    - starts-with: search for station name that starts with the provided string
  - distance: search by distance from a point
    - lat: latitude of the point
    - lon: longitude of the point
- Info
  - station-id: the station id
- Download
  - station-id: the station id
  - output: the file location to save the data
  - start: the start year
  - end: the end year

# TODO:
 - implement the download endpoint