package app

import (
	"endoscopy/internal/arge"
	"endoscopy/internal/logs"
	"endoscopy/internal/server"

	"github.com/urfave/cli/v2"
)

func New() (*cli.App, error) {
	apps := &cli.App{
		Commands: []*cli.Command{
			&cmdServer,
			&cmdVersion,
		},
	}
	return apps, nil
}

func runServer(c *cli.Context) error {
	err := arge.New(c)
	if err != nil {
		return err
	}
	logs.Info("Server Start  -- address: " + arge.Address + "  -- port: " + arge.Port)
	s := server.NewServer(arge.Address, arge.Port)
	s.RunServer()
	return nil
}
