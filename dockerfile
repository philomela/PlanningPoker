FROM mcr.microsoft.com/mssql/server

USER root

RUN apt-get -y update && \
		apt-get install -y golang-go

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY . /usr/src/app

RUN chmod +x /usr/src/app/Db/run-initialization.sh

EXPOSE 8080

USER mssql
ENTRYPOINT /bin/bash ./entrypoint.sh