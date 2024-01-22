#!/bin/bash
VERSION=v0.0.2

echo $VERSION

env GOOS=linux GOARCH=amd64 go build -o "bin/ggrep_${VERSION}_linux_x64" ./
env GOOS=linux GOARCH=arm64 go build -o bin/ggrep_${VERSION}_linux_arm64 ./
env GOOS=darwin GOARCH=arm64 go build -o bin/ggrep_${VERSION}_darwin_arm64 ./
env GOOS=darwin GOARCH=amd64 go build -o bin/ggrep_${VERSION}_darwin_x64 ./