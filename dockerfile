FROM golang:latest

ADD . /go/src/planning-poker
WORKDIR /go/src/planning-poker

RUN go build -o main .
ENTRYPOINT ["./main"]


