#!/bin/bash

echo "building..."

export GOPATH=$GOPATH:$(pwd)

#GOOS=darwin
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sendsms run/main.go

echo "build Success!!!"