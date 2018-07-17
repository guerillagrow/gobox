#!/bin/bash

#GOOS=linux GOARCH=arm GOARM=7  go build -o gobox_arm7 main.go

# build migdb_arm
# CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 go build -o migdb_arm migrate_db.go 


VERSION=`cat ./VERSION`

echo "Building GoBox"
echo "Compile GoBox ARM"
GOOS=linux GOARCH=arm GOARM=6 go build -ldflags "-X main.VERSION=${VERSION}" -v -o ./cmd/gobox/gobox_arm main.go
echo "Compile GoBox amd64"
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.VERSION=${VERSION}" -v -o ./cmd/gobox/gobox_amd64 main.go
echo "Compile GoBox x86"
GOOS=linux GOARCH=386 go build -ldflags "-X main.VERSION=${VERSION}" -v -o ./cmd/gobox/gobox_x86 main.go
echo "Building sensD"

cd ./cmd/sensd/
./build.sh
# ARM
#CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 go build -v -o gobox_arm7 main.go
# x86_64
#CC=x86_64-linux-gnu-gcc GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -v -o gobox_amd64 main.go

# playground
#CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 go build -v -o dev/playground_arm7 dev/playground.go
