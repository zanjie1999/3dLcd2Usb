del monitor-windows.exe
del monitor-linux
del monitor-darwin
SET CGO_ENABLED=1
SET GOARCH=
SET GOOS=windows
go build -ldflags "-H windowsgui"
rename monitor.exe monitor-windows.exe
SET CGO_ENABLED=0
SET GOARCH=amd64
SET GOOS=linux
go build
rename monitor monitor-linux
SET GOOS=darwin
SET GOARCH=amd64
go build
rename monitor monitor-darwin
