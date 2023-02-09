FROM mcr.microsoft.com/mssql/server

ENV ACCEPT_EULA=Y
ENV MSSQL_SA_PASSWORD=Yukon_900

EXPOSE 1433

USER root

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY . /usr/src/app

RUN chmod +x ./run-initialization.sh

USER mssql

CMD /bin/bash ./entrypoint.sh
