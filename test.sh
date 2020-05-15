#!/usr/bin/env sh

set -e

go test -v -count=1 -timeout=80s ./amqp
go test -v -count=1 -timeout=80s ./config
go test -v -count=1 -timeout=80s ./db
go test -v -count=1 -timeout=80s ./email
go test -v -count=1 -timeout=80s ./emit
go test -v -count=1 -timeout=80s ./http
go test -v -count=1 -timeout=80s ./log
go test -v -count=1 -timeout=80s ./redis
#go test -v -count=1 -timeout=80s ./tcp
# go test -v -count=1 -timeout=80 websocket