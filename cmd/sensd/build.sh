#!/bin/bash

echo "Compile sensD ARM"
CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 go build -v -o sensd_arm sensd.go

echo "Compile sensD amd64"
CC=x86_64-linux-gnu-gcc GOOS=linux GOARCH=amd64 CGO_ENABLED=1  go build -v -o sensd_amd64 sensd.go

echo "Compile sensD mock amd64"
GOARCH=amd64 GOOS=linux go build -v -o sensd_amd64_mock sensd_mock.go

echo "Compile sensD mock arm"
GOARCH=arm GOARM=6 GOOS=linux go build -v -o sensd_arm_mock sensd_mock.go