package app

import "github.com/urfave/cli/v2"

var (
	flagPath = cli.StringFlag{
		Name:  "path",
		Usage: "scan path",
	}
	flagLogPath = cli.StringFlag{
		Name:  "log",
		Value: "",
		Usage: "log file path",
	}
	flagOutput = cli.StringFlag{
		Name:  "output",
		Value: "",
		Usage: "输出文件目录",
	}
)

var cmdCli = cli.Command{
	Name:    "cli",
	Usage:   "命令行工具",
	Aliases: []string{"c"},
	Action:  runCli,
	Flags: []cli.Flag{
		&flagPath,
		&flagLogPath,
		&flagOutput,
	},
}
