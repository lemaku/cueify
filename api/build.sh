#/bin/bash
GOARCH=wasm GOOS=js go build -o main.wasm main.go
cp main.wasm ../main.wasm