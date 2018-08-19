FROM golang:alpine

RUN export GOBIN=/go/bin
RUN mkdir ${GOPATH}/src/app
COPY . ${GOPATH}/src/app
WORKDIR ${GOPATH}/src/app

RUN apk add --no-cache git
RUN go get && go build

ENTRYPOINT ${GOPATH}/src/app/app -bind=:8080