
echo Building GoBox
echo Compile GoBox ARM
GOOS=linux GOARCH=arm GOARM=6 go build -v -o cmd/gobox/gobox_arm main.go
echo Compile GoBox amd64
GOOS=linux GOARCH=amd64 go build -v -o cmd/gobox/gobox_amd64 main.go
echo Compile GoBox x86
GOOS=linux GOARCH=386 go build -v -o cmd/gobox/gobox_x86 main.go
echo Building sensD
