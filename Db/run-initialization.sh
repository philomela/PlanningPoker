#run the SQL_PlanningPoker script to create the DB and the schema in the DB
#do this in a loop because the timing for when the SQL instance is ready is indeterminate
for i in {1..50};
do
    /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P Yukon_900 -d master -i SQL_PlanningPoker.sql
    if [ $? -eq 0 ]
    then
        echo "SQL_PlanningPoker.sql completed"
        break
    else
        echo "not ready yet..."
        sleep 1
    fi
done
