#!/bin/bash
go build -ldflags="-s -w" main.go

chmod +x *.sh
