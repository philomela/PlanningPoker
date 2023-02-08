FROM golang:latest

RUN apt-get -y update && \
		apt-get install -y net-tools

ADD . /go/src/planning-poker
WORKDIR /go/src/planning-poker

RUN go build -o main .
ENTRYPOINT ["./main"]
