#!/bin/bash
FILE_EXTENSION=bak

export MSSQL_SA_PASSWORD=$DEFAULT_MSSQL_SA_PASSWORD
(/opt/mssql/bin/sqlservr --accept-eula & ) | grep -q "Server is listening on" && sleep 2

for restoreFile in /var/opt/mssql/data/*.$FILE_EXTENSION
do
    fileName=${restoreFile##*/}
    base=${fileName%.$FILE_EXTENSION}
    /opt/mssql-tools/bin/sqlcmd -S localhost -U SA -P $MSSQL_SA_PASSWORD -Q "RESTORE DATABASE [$base] FROM DISK = '$restoreFile'"
    rm -rf $restoreFile
done