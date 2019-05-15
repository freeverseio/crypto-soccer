package commands

import (
	"context"
	"sync"

	"github.com/urfave/cli"

	"github.com/freeverseio/go-soccer/config"
	"github.com/freeverseio/go-soccer/eth"
	"github.com/freeverseio/go-soccer/stakers"
	sto "github.com/freeverseio/go-soccer/storage"
	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	storage *sto.Storage
	stkrs   *stakers.Stakers
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

	return loadWeb3AndStakers(c)
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

func loadWeb3AndStakers(c *cli.Context) (err error) {

	loadConfig(c)

	// open rpc

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

	// open web3s
	var ks *keystore.KeyStore
	stks := make([]*stakers.Staker, 0, len(config.C.Stakers.Accounts))

	ks = keystore.NewKeyStore(config.C.Stakers.Keystore, keystore.StandardScryptN, keystore.StandardScryptP)
	for _, acc := range config.C.Stakers.Accounts {

		account, err := ks.Find(accounts.Account{
			Address: common.HexToAddress(acc.Address),
		})
		if err != nil {
			return err
		}

		err = ks.Unlock(account, config.C.Stakers.Keypasswd)
		if err != nil {
			return err
		}

		web3 := eth.NewWeb3Client(client, ks, &account)
		web3.ClientMutex = &sync.Mutex{}
		web3.MaxGasPrice = config.C.Web3.MaxGasPrice

		balance, err := client.BalanceAt(context.TODO(), account.Address, nil)
		if err != nil {
			return err
		}

		log.WithFields(log.Fields{
			"balance": balance.String(),
		}).Info(" Using account ", acc)

		stks = append(stks, &stakers.Staker{
			Address: account.Address,
			Client:  web3,
		})

	}
	stkrs, err = stakers.New(stks, storage)

	return err
}
