#!/usr/bin/env bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ./build/web-server main.go
#cp -rf ./configs ./build
#cp -rf ./deploy/* ./build/
mkdir ./build/website
cp -rf ./website/* ./build/website