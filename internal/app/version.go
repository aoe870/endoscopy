package app

import (
	"endoscopy/internal/version"

	"github.com/urfave/cli/v2"
)

var cmdVersion = cli.Command{
	Name:    "version",
	Usage:   "版本",
	Aliases: []string{"v"},
	Action:  version.RunVersion,
}
