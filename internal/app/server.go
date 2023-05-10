package app

import "github.com/urfave/cli/v2"

var (
	falgLog = cli.StringFlag{
		Name:  "log",
		Value: "analyzer.log",
		Usage: "日志文件，cli默认没日志",
	}

	flagPlatformHost = cli.StringFlag{
		Name:  "address",
		Value: "127.0.0.1",
		Usage: "平台回掉地址",
	}

	flagPort = cli.StringFlag{
		Name:  "port",
		Value: "9584",
		Usage: "服务端使用端口",
	}
)

var cmdServer = cli.Command{
	Name:    "server",
	Aliases: []string{"s"},
	Usage:   "服务器启动",
	Action:  runServer,
	Flags: []cli.Flag{
		&flagPlatformHost,
		&flagPort,
	},
}
