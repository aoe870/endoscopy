package arge

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"landau/internal/logs"
	"os"
)

var (
	Log string

	//server
	Port    string
	Address string
)

func New(c *cli.Context) error {

	Log = c.String("log")
	if Log != "" {
		logs.New(Log, true)
	} else {
		logs.New("analyzer.log", false)
	}

	//server
	Port = c.String("port")
	var env_ip string
	env_ip = os.Getenv("IP_SCA_PLATFORM")
	if len(env_ip) > 4 {
		Address = env_ip
	} else {
		Address = c.String("address")
	}
	if Address == "" {
		return errors.New("none address")
	}
	if Address == "" {
		return errors.New("none address")
	}

	//获取日志信息
	Log = c.String("log")

	return nil
}
