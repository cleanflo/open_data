#!/bin/bash

/opt/mssql-tools/bin/sqlcmd \
    -l 60 \
    -S localhost -U SA -P "$DEFAULT_MSSQL_SA_PASSWORD" \
    -Q "ALTER LOGIN SA WITH PASSWORD='${MSSQL_SA_PASSWORD}'" &

/opt/mssql/bin/permissions_check.sh "$@"