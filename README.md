# PlanningPoker.ru
This repo is site project PlanningPoker.ru
PlanningPoker - game for evaluating scrum team tasks

Run in Docker use command:
Step 1. docker build -t planning-poker .
Step 2. docker run --env-file env.list --name=planning-poker -p 80:8080 planning-poker