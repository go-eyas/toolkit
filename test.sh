#!/usr/bin/env sh

go test -v -timeout=80s ./amqp
go test -v -timeout=80s ./config
go test -v -timeout=80s ./http
go test -v -timeout=80s ./log
go test -v -timeout=80s ./redis
# go test -v -timeout=80 websocket