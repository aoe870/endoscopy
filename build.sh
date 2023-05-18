#!/usr/bin/env bash

outPut="Build"

#delete old build file
rm -rf outPut
mkdir outPut
cp README.md outPut

# Build the project to Windows
WindowsPath=outPut/windows
mkdir $WindowsPath
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -a \
                -gcflags="all=-trimpath=${PWD}" \
                -o $WindowsPath/endoscopy.exe cmd/main.go

# Build the project to Linux
LinuxPath=outPut/linux
mkdir $LinuxPath

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a \
                -gcflags="all=-trimpath=${PWD}" \
                -o $LinuxPath/endoscopy cmd/main.go

# Build the project to Mac
MacPath=outPut/mac
mkdir $MacPath
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -a \
                -gcflags="all=-trimpath=${PWD}" \
                -o $MacPath/endoscopy cmd/main.go
