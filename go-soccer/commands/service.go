package commands

import (
	"os"
	"os/signal"

	"github.com/freeverseio/go-soccer/service"
	"github.com/urfave/cli"
)

var ServiceCommands = []cli.Command{
	{
		Name:  "service",
		Usage: "manage services",
		Subcommands: []cli.Command{
			{
				Name:   "start",
				Usage:  "start the service",
				Action: serviceStart,
			},
		},
	},
}

func serviceStart(c *cli.Context) {

	must(load(c))
	serv, err := service.NewService(stkrs, storage)
	must(err)

	// catch ^C to send the stop signal
	ossig := make(chan os.Signal, 1)
	signal.Notify(ossig, os.Interrupt)

	go func() {
		for sig := range ossig {
			if sig == os.Interrupt {
				serv.Stop()
				os.Exit(1) // FIX
			}
		}
	}()

	serv.Start()
	serv.Join()
}
