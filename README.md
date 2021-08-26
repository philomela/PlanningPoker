# PlanningPoker.ru
This repo is site project PlanningPoker
PlanningPoker - game for evaluating scrum team tasks

Planning poker is a problem assessment technique based on team discussion of each problem. Each participant anonymously evaluates the problem by voting with a scorecard. Then the cards are revealed and the median for each problem is calculated.
#### Config preferences. 
For configuration project uses file ServersSettings.json and env.list if you would use docker container.
The settings files contain values ​​* instead of which you need to substitute your own data this is done for security reasons
#### Restore database.
For working this project, restoring database and start run sql scripts from file "SQL_PlanningPoker.sql".
* Step 1. Open file with sql scripts. Highlights all scripts and press F5 or run with MSSQL Managment Studio.
* Step 2. Wait finish work. 
#### Run in Docker use command:
* Step 1. ```docker build -t planning-poker .```
* Step 2. ```docker run --env-file env.list --name=planning-poker -p 80:8080 planning-poker```

Run Project in VSCode or anything:
For run project enter command in terminal near project "go run main.go" or "go build main.go" and current directory run main executable file.

enjoy!:heart:
