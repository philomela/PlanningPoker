FROM mcr.microsoft.com/mssql/server

ENV ACCEPT_EULA=Y
ENV MSSQL_SA_PASSWORD=Yukon_900

USER root

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY . /usr/src/app

RUN chmod +x ./run-initialization.sh

EXPOSE 1433

USER mssql

CMD /bin/bash ./entrypoint.sh
