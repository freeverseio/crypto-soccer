package commands

import (
	"math/big"

	"github.com/freeverseio/go-soccer/lionel"
	sto "github.com/freeverseio/go-soccer/storage"
	"github.com/urfave/cli"
)

var TestCommands = []cli.Command{
	{
		Name:   "test",
		Usage:  "run test",
		Action: runTest,
	},
}

func runTest(c *cli.Context) {

	userActions := []sto.UserActions{
		// day 0
		sto.UserActions{
			Tactics: [][3]uint8{[3]uint8{4, 4, 2}, [3]uint8{4, 3, 3}},
		},
		// day 1
		//		sto.UserActions{
		//			Tactics: [][3]uint8{[3]uint8{4, 3, 3}, [3]uint8{4, 4, 3}},
		//		},
	}

	must(load(c))
	lionel, err := lionel.New(stkrs.Members()[0].Client, storage, stkrs)
	must(err)
	_, err = lionel.ComputeLeague(
		big.NewInt(0),
		[]*big.Int{big.NewInt(1), big.NewInt(2)},
		userActions,
		false,
	)
	must(err)

}
