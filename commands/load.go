package commands

import (
	"context"
	"sync"

	cfg "github.com/freeverseio/go-soccer/config"
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

func load() error {

	var err error

	if err = loadStorage(); err != nil {
		return err
	}

	return loadWeb3()
}

func loadStorage() error {

	var err error

	storage, err = sto.New(cfg.C.DB.Path)

	return err

}

func loadWeb3() (err error) {

	// load keystore

	var ks *keystore.KeyStore
	var account accounts.Account

	ks = keystore.NewKeyStore(cfg.C.Keystore.Path, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err = ks.Find(accounts.Account{
		Address: common.HexToAddress(cfg.C.Keystore.Account),
	})
	if err != nil {
		return err
	}

	err = ks.Unlock(account, cfg.C.Keystore.Passwd)
	if err != nil {
		return err
	}

	log.WithField("acc", cfg.C.Keystore.Account).Info("Account unlocked")

	// load web3
	log.WithField("url", cfg.C.Web3.RPCURL).Info("Checking WEB3.")

	client, err := ethclient.Dial(cfg.C.Web3.RPCURL)
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

	return nil
}
