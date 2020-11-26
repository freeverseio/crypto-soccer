package consumer

import (
	"crypto/ecdsa"

	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionpassmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/googleplaystoreutils"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	log "github.com/sirupsen/logrus"
)

func ProcessAuctionPassPlaystoreOrders(
	service storage.Tx,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	googleCredentials []byte,
	namesdb *names.Generator,
	iapTestOn bool,
) error {
	orders, err := service.AuctionPassPlayStorePendingOrders()
	if err != nil {
		return err
	}

	client, err := googleplaystoreutils.NewGoogleClientService(googleCredentials)
	if err != nil {
		return err
	}

	for _, order := range orders {
		machine, err := auctionpassmachine.New(
			service,
			client,
			order,
			contracts,
			pvc,
			iapTestOn,
		)
		if err != nil {
			return err
		}
		if err := machine.Process(); err != nil {
			log.Error(err)
			continue
		}
		if err := service.AuctionPassPlayStoreUpdateState(machine.Order()); err != nil {
			return err
		}
	}

	return nil
}
