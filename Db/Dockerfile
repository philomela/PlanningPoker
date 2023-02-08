FROM mcr.microsoft.com/mssql/server:2019-latest
FROM golang:latest

ENV ID=sa
ENV ACCEPT_EULA=Y
ENV MSSQL_SA_PASSWORD=Yukon900
ENV MSSQL_PID=Express

EXPOSE 1433


RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY . /usr/src/app

RUN go build -o main .

USER root

RUN chmod +x /usr/src/app/Db/run-initialization.sh

USER mssql

CMD /bin/bash Db/entrypoint.sh
ENTRYPOINT ["main"]