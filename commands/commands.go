package commands

import (
	"errors"
	"os"
	"os/signal"

	"github.com/freeverseio/go-soccer/service"
	sto "github.com/freeverseio/go-soccer/storage"

	"github.com/spf13/cobra"
)

var (
	errInvalidParameters = errors.New("invalid parameters")
)

func CmdDumpDb(cmd *cobra.Command, args []string) {

	must(loadStorage())

	storage.Dump(os.Stdout)
}

func CmdInitDb(cmd *cobra.Command, args []string) {

	must(loadStorage())

	storage.SetGlobals(sto.GlobalsEntry{})

}

func CmdServe(cmd *cobra.Command, args []string) {

	must(load())
	serv, err := service.NewService(web3, storage)
	must(err)

	// catch ^C to send the stop signal
	ossig := make(chan os.Signal, 1)
	signal.Notify(ossig, os.Interrupt)
	go func() {
		for sig := range ossig {
			if sig == os.Interrupt {
				serv.Stop()
			}
		}
	}()
	serv.Join()
}
