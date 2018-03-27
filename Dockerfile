FROM golang:alpine

WORKDIR /go/src/github.com/jeanlucthumm/thummcoin

RUN apk add --no-cache bash

ENTRYPOINT ["./devrun.sh"]