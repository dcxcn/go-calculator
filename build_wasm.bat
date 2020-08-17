cd "%~dp0"
set GOARCH=wasm
set GOOS=js
go build -o wasm/go-calculator.wasm main.go