# PlanningPoker
PlanningPoker - game for evaluating scrum team tasks

Planning poker is a problem assessment technique based on team discussion of each problem. Each participant anonymously evaluates the problem by voting with a scorecard. Then the cards are revealed and the median for each problem is calculated.
#### Config preferences. 
For configuration project uses file ServersSettings.json and env.list if you would use docker container.
The settings files contain values ​​* instead of which you need to substitute your own data this is done for security reasons
#### Restore local database for local start app without Docker.
For local working this project, restoring database and start run sql scripts from file "SQL_PlanningPoker.sql".
* Step 1. Open file with sql scripts. Highlights all scripts and press F5 or run with MSSQL Managment Studio.
* Step 2. Wait finish work. 

#### Run in local use command:
* Step 1. ```cd <YourPathToProject>```
* Step 2. Change server in connection string on localhost 
```server=localhost;user id=PlanningPoker;password=Pa$$word;database=PlanningPoker;```
* Step 3. ```go run main.go```
* Step 4. ```localhost``` - address app.

#### Run in Docker use command:
* Step 1. ```cd <YourPathToProject>```
* Step 2. Change server in connection string on mssql 
```server=mssql;user id=PlanningPoker;password=Pa$$word;database=PlanningPoker;```
* Step 3. ```docker-compose up -d --build --force-recreate```
* Step 4. Wait db inti scripts and ```localhost``` address app

enjoy!:heart:
