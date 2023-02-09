#start SQL Server, start the script to create the DB and import the data, start the app
/usr/src/app/Db/run-initialization.sh & /opt/mssql/bin/sqlservr & go run /usr/src/app/main.go