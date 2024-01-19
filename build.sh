#!/bin/bash

env GOOS=linux GOARCH=amd64 go build -o bin/ggrep_v0.0.1_linux_x64 ./
env GOOS=linux GOARCH=arm64 go build -o bin/ggrep_v0.0.1_linux_arm64 ./
env GOOS=darwin GOARCH=arm64 go build -o bin/ggrep_v0.0.1_mac_arm64 ./
env GOOS=darwin GOARCH=amd64 go build -o bin/ggrep_v0.0.1_mac_x64 ./
