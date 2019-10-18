package commands

import (
	"os"

	sto "github.com/freeverseio/go-soccer/storage"
	"github.com/urfave/cli"
)

var DbCommands = []cli.Command{
	{
		Name:  "db",
		Usage: "create and manage db",
		Subcommands: []cli.Command{
			{
				Name:   "init",
				Usage:  "Initialize the database",
				Action: DbInit,
			},
			{
				Name:   "dump",
				Usage:  "DbInit the database",
				Action: DbDump,
			}},
	},
}

func DbDump(c *cli.Context) {

	must(loadStorage(c))

	storage.Dump(os.Stdout)
}

func DbInit(c *cli.Context) {

	must(loadStorage(c))

	storage.SetGlobals(sto.GlobalsEntry{})

}
