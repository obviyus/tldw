package commands

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"tldw-server/internal/config"
	"tldw-server/internal/server"
)

var startFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:   "detach-server, d",
		Usage:  "detach from the console (daemon mode)",
		EnvVar: "DOWNCOUNT_DETACH_SERVER",
	},
	&cli.BoolFlag{
		Name:  "config, c",
		Usage: "show config",
	},
}

var StartCommand = cli.Command{
	Name:    "start",
	Aliases: []string{"up"},
	Usage:   "Starts web server",
	Flags:   startFlags,
	Action:  startAction,
}

func startAction(ctx *cli.Context) error {
	conf := config.NewConfig(ctx)

	if ctx.IsSet("config") {
		return nil
	}

	// pass this context down the chain
	cctx, cancel := context.WithCancel(context.Background())

	if err := conf.Init(); err != nil {
		log.Fatalln(err)
	}

	// initialize the database
	conf.InitDb()

	// start web server
	go server.Start(cctx)

	// set up proper shutdown of daemon and web server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	logrus.Info("shutting down...")
	cancel()

	return nil
}
