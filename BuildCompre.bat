set GOOS=windows
packr
go generate
go build -ldflags="-s -w"
packr clean
del /f PokecatcherCompre.exe
D:\upx\upx.exe -9 -oPokecatcherCompre.exe "%~dp0Pokecatcher.exe"
PAUSE
PAUSE