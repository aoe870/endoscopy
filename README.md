<br><br>

<h1 align="center">Endoscopy</h1>

<p align="center">
    <a href="/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg"/></a>
    <img alt="Github Test" src="https://github.com/zricethezav/gitleaks/actions/workflows/test.yml/badge.svg">  
    <a href="https://goreportcard.com/report/github.com/mingrammer/commonregex"><img src="https://goreportcard.com/badge/github.com/mingrammer/commonregex"/></a>
    <img src="https://img.shields.io/twitter/follow/zricethezav?label=Follow%20zricethezav&style=social&color=blue" alt="Follow @taotao01114978" />
</p>
<p align="center">
  A tool written in go to detect sensitive information
</p>

## Introduction
endoscopy is intended for scanning third-party libraries for sensitive information. It is capable of finding secrets accidentally committed to a git repo, additional credentials provided along with compromised credentials and secrets stored in plaintext/config files. The goal of this tool is to increase awareness regarding the types of sensitive information that are often accidentally shared.

## Usage
```
aoe@computer cmd % ./endoscopy -h 
NAME:
   endoscopy - A new cli application

USAGE:
   endoscopy [global options] command [command options] [arguments...]

COMMANDS:
   server, s   服务器启动
   version, v  版本
   cli, c      命令行工具
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help

```
if you want to use cli, you can use this command
```
aoe@computer cmd %./endoscopy c -h
NAME:
   endoscopy cli - 命令行工具

USAGE:
   endoscopy cli [command options] [arguments...]

OPTIONS:
   --path value    scan path
   --log value     log file path
   --output value  输出文件目录
   --help, -h      show help
```
## Installing

```
go get github.com/aoe870/endoscopy.git
```

if you want build binary file, you can use this command
###  use build shell script
``` 
./build.sh
```
###  use go build
```
go build -o endoscopy cmd/main.go
```