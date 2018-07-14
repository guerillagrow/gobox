#!/bin/bash

RPI_IP=$1

echo "Upload GoBox to your Raspberry Pi"
echo "Raspberry Pi IP-Address: ${RPI_IP}"
echo "-----------------------------------"


echo "Upload GoBox service (init.d)"
scp -r ./initscript/gobox.sh root@$RPI_IP:/etc/init.d/gobox

echo "Upload GoBox config"
scp -r ./conf/ root@$RPI_IP:/usr/local/gobox/ 

echo "Upload GoBox views"
scp -r ./views/ root@$RPI_IP:/usr/local/gobox/ 

echo "Upload GoBox static files"
scp -r ./static/ root@$RPI_IP:/usr/local/gobox/ 

echo "Upload GoBox export folder"
scp -r ./export/ root@$RPI_IP:/usr/local/gobox/

echo "Upload GoBox tmp folder"
scp -r ./tmp/ root@$RPI_IP:/usr/local/gobox/ 

echo "Upload GoBox binary"
scp ./cmd/gobox/gobox_arm root@$RPI_IP:/usr/local/bin/gobox

echo "Upload GoBox sensD binary"
scp ./cmd/sensd/sensd_arm root@$RPI_IP:/usr/local/bin/sensd
