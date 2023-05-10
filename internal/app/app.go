package app

import (
	"github.com/urfave/cli/v2"
	"landau/internal/arge"
	"landau/internal/logs"
	"landau/internal/server"
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
	arge.New(c)
	logs.Info("Server Start  -- address: " + arge.Address + "  -- port: " + arge.Port)
	s := server.NewServer(arge.Address, arge.Port)
	s.RunServer()
	return nil
}
