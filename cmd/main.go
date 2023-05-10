package main

import (
	"landau/internal/app"
	"landau/internal/logs"
	"os"
)

func main() {
	app, err := app.New()
	if err != nil {
		logs.Error(err)
	}
	err = app.Run(os.Args)
	if err != nil {
		logs.Error(err)
	}
}
