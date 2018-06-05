FROM golang:alpine

WORKDIR /go/src/github.com/AlexsJones/gravitywell
COPY . .

RUN apk --no-cache add git && \
    set -x && \
    go get -v -d ./... && \
    go install && \
    go test --cover -v ./... && \
    rm -rf /go/src /go/pkg

WORKDIR /home/root
