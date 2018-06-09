#!/bin/bash

mkdir -p build
GOOS=linux GOARCH=amd64 go build -o build/scm-status-linux-amd64
docker build -t jimmysawczuk/scm-status .
docker push jimmysawczuk/scm-status
rm -rf build
