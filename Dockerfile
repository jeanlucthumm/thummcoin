FROM golang:alpine

WORKDIR /go/src/github.com/jeanlucthumm/thummcoin

RUN apk update && apk add bash git

ENTRYPOINT ["./devrun.sh"]