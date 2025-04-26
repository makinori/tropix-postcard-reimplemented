#!/bin/sh
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o Postcard.exe main.go