FROM golang:alpine

WORKDIR /go/src/app

RUN apk add --no-cache bash

ENTRYPOINT ["./devrun.sh"]