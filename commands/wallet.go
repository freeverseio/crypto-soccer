package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/freeverseio/go-soccer/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var WalletCommands = []cli.Command{
	{
		Name:  "wallet",
		Usage: "create and manage eth wallet",
		Subcommands: []cli.Command{
			{
				Name:   "new",
				Usage:  "create new Keystorage. Args: keystoreId password",
				Action: cmdWalletNew,
			},
			{
				Name:   "info",
				Usage:  "display info about the current relay wallets",
				Action: cmdWalletInfo,
			}},
	},
}

func cmdWalletNew(c *cli.Context) error {
	loadConfig(c)

	if config.C.Keystore.Passwd == "" {
		log.Fatal("No password specified.")
	}

	ks := keystore.NewKeyStore(config.C.Keystore.Path, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(config.C.Keystore.Passwd)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("New Keystore with address '", account.Address.Hex()+"' created in '"+config.C.Keystore.Path+"'")
	return nil
}

func cmdWalletInfo(c *cli.Context) error {
	loadConfig(c)

	files, err := ioutil.ReadDir(config.C.Keystore.Path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("current Keystores:")
	for _, f := range files {
		fmt.Println("	", f.Name())
	}
	return nil
}
