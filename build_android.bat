set GOOS=android
set GOARCH=arm64
set CGO_ENABLED=0
go build -ldflags="-s -w"