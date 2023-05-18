package arge

import (
	"endoscopy/internal/logs"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
)

var (
	Log string

	// server
	Port    string
	Address string

	// cli
	Path   string
	Output string
)

func New(c *cli.Context) error {

	Log = c.String("log")
	if Log != "" {
		logs.New(Log, true)
	} else {
		logs.New("endoscopy.log", false)
	}

	// cli
	if c.Command.Name == "cli" {
		if c.String("path") != "" {
			dir, filename := filepath.Split(c.String("path"))
			if len(filename) > 0 {
				Path = c.String("path")
			} else {
				if dir == "/" {
					Path = dir
				} else {
					Path = dir[0 : len(dir)-1]
				}
			}
		}
		Output = c.String("output")
	} else {
		//server
		Address = os.Getenv("IP_SCA_PLATFORM")
	}
	//获取日志信息
	Log = c.String("log")

	return nil
}
