#!/bin/bash

RPI_IP=$1
##echo "PI IP=${RPI_IP}"

echo "Upload GoBox service (init.d)"
scp -r ./initscript/gobox.sh root@$RPI_IP:/etc/init.d/gobox

echo "Upload GoBox config"
scp -r ./conf/ root@$RPI_IP:/usr/local/gobox/ 

echo "Upload GoBox views"
scp -r ./views/ root@$RPI_IP:/usr/local/gobox/ 

echo "Upload GoBox static files"
scp -r ./static/ root@$RPI_IP:/usr/local/gobox/ 

echo "Upload GoBox binary"
scp ./cmd/gobox/gobox_arm root@$RPI_IP:/usr/local/bin/gobox

echo "Upload GoBox sensD binary"
scp ./cmd/sensd/sensd_arm root@$RPI_IP:/usr/local/bin/sensd
