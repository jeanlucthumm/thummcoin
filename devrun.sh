#!/bin/bash

go get -d -v .
go install -v .

exec "$GOPATH/bin/thummcoin" "$@"