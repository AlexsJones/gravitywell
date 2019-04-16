FROM golang:latest as Builder

WORKDIR /go/src/github.com/AlexsJones/gravitywell
COPY . .

RUN set -x && \
    export GOPATH="/go" && export PATH=$PATH:$GOPATH/bin && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && \
    dep ensure -v && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X 'main.version=$(cat VERSION)' -X 'main.revision=$(git rev-parse --short HEAD)' -X 'main.buildtime=$(date -u +%Y-%m-%d.%H:%M:%S)'" -o /gravitywell && \
    rm -rf ${GOPATH}

FROM alpine:3.8

COPY --from=Builder /gravitywell /usr/bin/gravitywell

RUN apk --no-cache add dumb-init

ENTRYPOINT ["dumb-init", "gravitywell"]
