# PlanningPoker.ru
This repo is site project PlanningPoker.ru
PlanningPoker - game for evaluating scrum team tasks

Restore database.
For working this project, restoring database and start run sql scripts from file "SQL_PlanningPoker.sql".
    1. Open file with sql scripts. Highlights all scripts and press F5 or run with MSSQL Managment Studio.
    2. Wait finish work. 

Config preferences. 
For configuration project uses file ServersSettings.json and env.list if you would use docker container.

Run in Docker use command:
Step 1. docker build -t planning-poker .
Step 2. docker run --env-file env.list --name=planning-poker -p 80:8080 planning-poker

Run Project in VSCode or anything:
For run project enter command in terminal near project "go run main.go" or "go build main.go" and current directory run main executable file.