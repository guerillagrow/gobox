echo off

set rpiip=%1

echo "Upload GoBox to your Raypberry Pi"
echo "Raspberry PI IP-Address: %rpiip%"
echo "-----------------------------------"

echo "Upload GoBox service (init.d)"
pscp -r ./initscript/gobox.sh root@%rpiip%:/etc/init.d/gobox

echo "Upload GoBox config"
pscp -r ./conf/ root@%rpiip%:/usr/local/gobox/ 

echo "Upload GoBox views"
pscp -r ./views/ root@%rpiip%:/usr/local/gobox/ 

echo "Upload GoBox static files"
pscp -r ./static/ root@%rpiip%:/usr/local/gobox/ 

echo "Upload GoBox export folder"
pscp -r ./export/ root@%rpiip%:/usr/local/gobox/

echo "Upload GoBox tmp folder"
pscp -r ./tmp/ root@%rpiip%:/usr/local/gobox/ 

echo "Upload GoBox binary"
pscp ./cmd/gobox/gobox_arm root@%rpiip%:/usr/local/bin/gobox

echo "Upload GoBox sensD binary"
pscp ./cmd/sensd/sensd_arm root@%rpiip%:/usr/local/bin/sensd
