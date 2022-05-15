# Water Well Data

This package provides a simple API for accessing published water well data from Canadian provinces. The published data was collected in a variety of formats(CSV & MDB mainly) and imported into MSSQL databases, some of the databases have multiple tables. These databases were then exported for building Docker containers to simplify the hosting requirements, the build files for the Docker images can be found in the src directory, the actual source files and the exported MDB files are available at these links:

### Data Sources
[sources.tar.gz](https://github.com/cleanflo/open_data/releases/download/src-0.0.1/sources.tar.gz)

- [Government of Alberta](http://groundwater.alberta.ca/WaterWells/d/#tabs_tablist_dijit_layout_ContentPane_5)
 - [Data Model](http://groundwater.alberta.ca/WaterWells/docs/Awwid%20Data%20Model.pdf)
- [Government of British Columbia](https://apps.nrs.gov.bc.ca/gwells/)
 - [Licence](https://www2.gov.bc.ca/gov/content?id=A519A56BC2BF44E4A008B33FCF527F61)
- [Government of Ontario](https://data.ontario.ca/dataset/well-records)
 - [Licence](https://www.ontario.ca/page/open-government-licence-ontario)
- [Government of Nova Scotia](https://www.novascotia.ca/nse/groundwater/welldatabase.asp)
 - [User Guide](https://www.novascotia.ca/nse/groundwater/docs/UsersManual_NSWellLogsDatabase.pdf)
- [Government of Saskatchewan](https://gis.wsask.ca/Html5Viewer/index.html?viewer=WaterWells.WellsViewer/)

Datasets Not Included:
- [Government of Quebec](https://www.environnement.gouv.qc.ca/eau/souterraines/sih/index.htm)
  - [Licence](https://www.donneesquebec.ca/fr/licence/#cc-by)
  - Reason: Data is only published in ESRI

### Exported from MSSQL
The MDBs were used in the creation of the docker images.
[data.tar.gz](https://github.com/cleanflo/open_data/releases/download/src-0.0.1/data.tar.gz)

The docker images can be built using the docker build command. The recommended process is to clone this repository, download the `data.tar.gz` file into the `src` directory, and run docker build command targetting the desired image.
```
git clone git@github.com:cleanflo/open_data.git
cd open_data/
wget https://github.com/cleanflo/open_data/releases/download/src-0.0.1/data.tar.gz
docker build --tag=cleanflo/well-data_all-mssql --target=all-data .
```

The docker images are published on Docker Hub, the links to the images are:
[Docker Hub Images](https://hub.docker.com/repository/docker/cleanflo/well-data_all-mssql)


## Usage
The current means of using the API is via the `Handler` method, which is a `net/http.HandlerFunc` that handles the request, performs the query, and returns the data. The `Handler` method can be used as below, where `well-data` is the name of the docker container that contains the MSSQL database:

```go
wellsRouter := dataRouter.PathPrefix("/wells/").Subrouter()

wellsRouter.HandleFunc("/alberta/", wells.AlbertaR.User("sa").Password($MSSQL_PASSWORD).Host("well-data").Port(1433).Handler).Methods("GET")

wellsRouter.HandleFunc("/british-columbia/", wells.BritishColumbiaR.User("sa").Password($MSSQL_PASSWORD).Host("well-data").Port(1433).Handler).Methods("GET")

wellsRouter.HandleFunc("/nova-scotia/", wells.NovaScotiaR.User("sa").Password($MSSQL_PASSWORD).Host("well-data").Port(1433).Handler).Methods("GET")

wellsRouter.HandleFunc("/ontario/", wells.OntarioR.User("sa").Password($MSSQL_PASSWORD).Host("well-data").Port(1433).Handler).Methods("GET")

wellsRouter.HandleFunc("/saskatchewan/", wells.SaskatchewanR.User("sa").Password($MSSQL_PASSWORD).Host("well-data").Port(1433).Handler).Methods("GET")

```

## Description
This package attempts to abstract the query logic and provide a simple http API for querying the databases by building predictable SQL statements. The requests are made using GET parameters to enable response caching on a variety of host providers, and the response is chunked and returned as a JSON. The JSON indicates the number of records returned, the total number of records that satisfy the query, the total number of pages, and the current page. The response data is currently limited to the latitude and longitude of the well, but PRs and suggestions are welcome. As our current usecase requires only the two fields, the entire struture of this project is open for revisioning and extension.

The API for each endpoint can perform filtering on certain parameters, but depending on the published data the mapping may be incomplete, untracked, or excluded for a particular province. Many of the columns of interest were counted and identified each parameter to perform mapping. Review each province's source file for a list of the "columns of interest" and the mapping to the API.

### Filter by Columns
- Completed[start:end]
 - start - time.Time
 - end - time.Time
- Abandoned[start:end]
 - start - time.Time
 - end - time.Time
- Status[1,2,3..,n] - string
- Use[1,2,3..,n] - string
- Colour[1,2,3..,n] - string
- Taste[1,2,3..,n] - string
- Odour[1,2,3..,n] - string
- Rate[start:end]
 - start - int
 - end - int
- Depth[start:end]
 - start - int
 - end - int
- Bedrock[start:end]
 - start - int
 - end - int

# TODO:
 - Work on docker build process to support other databases
 - Extract the `ClientRequestData` struct from the `db_layer` subpackage and create an means of providing the request data to the `Retreiver` struct
 - Create an endpoint for returning which fields and filters are available for a defined `Retreiver` 
 - Enhance the README with a full description of each provinces data and API description.
 - DB layer subpackage should get better testing
 - Logging lines throughout these packages should be reviewed and cleaned