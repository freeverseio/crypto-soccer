package consumer

import (
	"crypto/ecdsa"
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/notary/playstore"
	log "github.com/sirupsen/logrus"
)

func ProcessPlaystoreOrders(
	tx *sql.Tx,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	googleCredentials []byte,
	iapTestOn bool,
) error {
	service := postgres.NewPlaystoreOrderHistoryService(tx)

	orders, err := service.PendingOrders()
	if err != nil {
		return err
	}

	client, err := playstore.NewGoogleClientService(googleCredentials)
	if err != nil {
		return err
	}

	for _, order := range orders {
		machine, err := playstore.New(
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
		if err := service.UpdateState(machine.Order()); err != nil {
			return err
		}
	}

	return nil
}
