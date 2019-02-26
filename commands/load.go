package commands

import (
	"context"
	"sync"

	"github.com/urfave/cli"

	"github.com/freeverseio/go-soccer/config"
	eth "github.com/freeverseio/go-soccer/eth"
	sto "github.com/freeverseio/go-soccer/storage"
	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	web3    *eth.Web3Client
	storage *sto.Storage
)

func must(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func load(c *cli.Context) error {

	var err error

	if err = loadStorage(c); err != nil {
		return err
	}

	return loadWeb3(c)
}

func loadConfig(c *cli.Context) {
	var err error
	if err = config.MustRead(c); err != nil {
		log.Fatal(err)
	}
}

func loadStorage(c *cli.Context) error {

	loadConfig(c)

	var err error
	storage, err = sto.New(config.C.DB.Path)

	return err

}

func loadWeb3(c *cli.Context) (err error) {

	loadConfig(c)

	// open wallet
	var ks *keystore.KeyStore
	var account accounts.Account

	ks = keystore.NewKeyStore(config.C.Keystore.Path, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err = ks.Find(accounts.Account{
		Address: common.HexToAddress(config.C.Keystore.Account),
	})
	if err != nil {
		return err
	}

	err = ks.Unlock(account, config.C.Keystore.Passwd)
	if err != nil {
		return err
	}

	// load web3

	log.WithField("url", config.C.Web3.RPCURL).Info("Checking WEB3.")

	client, err := ethclient.Dial(config.C.Web3.RPCURL)
	if err != nil {
		return err
	}

	clientnetworkid, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}
	log.WithField("networkid", clientnetworkid).Info("Checking WEB3.")

	web3 = eth.NewWeb3Client(client, ks, &account)
	web3.ClientMutex = &sync.Mutex{}
	web3.MaxGasPrice = config.C.Web3.MaxGasPrice

	balance, err := client.BalanceAt(context.TODO(), account.Address, nil)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"balance": balance.String(),
	}).Info(" Using account ", config.C.Keystore.Account)

	return nil
}
