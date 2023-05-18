<br><br>

<h1 align="center">Endoscopy</h1>

<p align="center">
  <a href="/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg"/></a>
  <a href="https://app.fossa.io/projects/git%2Bgithub.com%2Fmingrammer%2Fcommonregex?ref=badge_shield" alt="FOSSA Status"><img src="https://app.fossa.io/api/projects/git%2Bgithub.com%2Fmingrammer%2Fcommonregex.svg?type=shield"/></a>
  <a href="https://godoc.org/github.com/mingrammer/commonregex"><img src="https://godoc.org/github.com/mingrammer/commonregex?status.svg"/></a>
  <a href="https://goreportcard.com/report/github.com/mingrammer/commonregex"><img src="https://goreportcard.com/badge/github.com/mingrammer/commonregex"/></a>
  <a href="https://codecov.io/gh/mingrammer/commonregex"><img src="https://codecov.io/gh/mingrammer/commonregex/branch/master/graph/badge.svg" /></a>
</p>
<p align="center">
  A tool written in go to detect sensitive information
</p>

<br>

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
#  use build shell script
``` 
./build.sh
```

#  use go build
```
go build -o endoscopy cmd/main.go
```