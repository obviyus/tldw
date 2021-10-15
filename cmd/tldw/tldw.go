package main

import (
	"os"

	"github.com/urfave/cli"

	"tldw-server/internal/commands"
	"tldw-server/internal/event"
)

var version = "alpha"
var log = event.Log

func main() {
	app := cli.NewApp()

	app.Name = "TL;DW"
	app.Usage = "Server to run the API for TL;DW"
	app.Version = version
	app.Copyright = "(c) 2021-2021 Ayaan Zaidi <hi@obviy.us>"
	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		commands.StartCommand,
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}
