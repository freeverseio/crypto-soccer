package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/freeverseio/go-soccer/commands"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gosoccer"
	app.Version = "0.0.1-alpha"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config"},
	}

	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, commands.WalletCommands...)
	app.Commands = append(app.Commands, commands.DbCommands...)
	app.Commands = append(app.Commands, commands.ServiceCommands...)
	err := app.Run(os.Args)
	if err != nil {
		color.Red(err.Error())
	}
}
