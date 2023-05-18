package app

import (
	"endoscopy/internal/arge"
	"endoscopy/internal/logs"
	"endoscopy/internal/server"
	"endoscopy/internal/task"
	"math/rand"
	"time"

	"github.com/urfave/cli/v2"
)

func New() (*cli.App, error) {
	apps := &cli.App{
		Commands: []*cli.Command{
			&cmdServer,
			&cmdVersion,
			&cmdCli,
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

func runCli(c *cli.Context) error {
	err := arge.New(c)
	if err != nil {
		return err
	}

	taskConfig := task.Config{
		OutPutPath: arge.Output,
		InPutPath:  arge.Path,
	}

	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 20; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	taskConfig.TaskID = string(result)
	task.New(taskConfig)
	return nil
}
