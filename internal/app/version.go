package app

import (
	"github.com/urfave/cli/v2"
	"landau/internal/version"
)

var cmdVersion = cli.Command{
	Name:    "version",
	Usage:   "版本",
	Aliases: []string{"v"},
	Action:  version.RunVserion,
}
