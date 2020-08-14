cd "%~dp0"
call make_version.bat ./ version.h

cd "%~dp0"
windres.exe -i main.rc -o main.syso

call go build -ldflags "-H windowsgui -w -s"