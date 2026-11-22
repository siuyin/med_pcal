#!/bin/sh
GOOS=android GOARCH=arm64 CGO_ENABLED=0 go build -o calSvr ./cmd/calSvr
