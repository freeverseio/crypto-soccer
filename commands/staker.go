package commands

import (
	"crypto/rand"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/go-soccer/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var StakerCommands = []cli.Command{
	{
		Name:  "staker",
		Usage: "manage PoA staker",
		Subcommands: []cli.Command{
			{
				Name:   "info",
				Usage:  "get staker information",
				Action: stakerInfo,
			},
			{
				Name:   "new",
				Usage:  "create new staker account",
				Action: stakerNew,
			},
			{
				Name:   "enroll",
				Usage:  "enroll as a staker",
				Action: stakerEnroll,
			},
			{
				Name:   "queryunenroll",
				Usage:  "query to unenroll a staker",
				Action: stakerQueryUnenroll,
			},
			{
				Name:   "unenroll",
				Usage:  "unenroll a staker",
				Action: stakerUnenroll,
			},
		},
	},
}

func stakerNew(c *cli.Context) error {
	loadConfig(c)

	if config.C.Stakers.Keypasswd == "" {
		log.Fatal("No password specified.")
	}

	ks := keystore.NewKeyStore(config.C.Stakers.Keystore, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(config.C.Stakers.Keypasswd)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("New Keystore with address ", account.Address.Hex()+" created in '"+config.C.Stakers.Keystore+"'")
	var onionhash common.Hash
	if _, err := rand.Read(onionhash[:]); err != nil {
		return err
	}
	log.Info("Onionhash '", onionhash.Hex())
	return nil
}

func stakerInfo(c *cli.Context) error {
	return nil
}

func stakerEnroll(c *cli.Context) error {
	loadConfig(c)
	if err := loadWeb3AndStakers(c); err != nil {
		return err
	}
	if err := load(c); err != nil {
		return err
	}
	if len(c.Args().Get(0)) == 0 {
		return fmt.Errorf("Needs staker address to enroll")
	}
	address := common.HexToAddress(c.Args().Get(0))
	if err := stkrs.Enroll(address); err != nil {
		return err
	}
	return nil
}

func stakerQueryUnenroll(c *cli.Context) error {
	loadConfig(c)
	if err := loadWeb3AndStakers(c); err != nil {
		return err
	}
	if err := load(c); err != nil {
		return err
	}
	if len(c.Args().Get(0)) == 0 {
		return fmt.Errorf("Needs staker address to query un-enroll")
	}
	address := common.HexToAddress(c.Args().Get(0))
	if err := stkrs.QueryUnenroll(address); err != nil {
		return err
	}
	return nil
}
func stakerUnenroll(c *cli.Context) error {
	loadConfig(c)
	if err := loadWeb3AndStakers(c); err != nil {
		return err
	}
	if err := load(c); err != nil {
		return err
	}
	if len(c.Args().Get(0)) == 0 {
		return fmt.Errorf("Needs staker address to un-enroll")
	}
	address := common.HexToAddress(c.Args().Get(0))
	if err := stkrs.Unenroll(address); err != nil {
		return err
	}
	return nil
}
