package consumer

import (
	"crypto/ecdsa"
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/notary/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func ProcessPlaystoreOrders(
	tx *sql.Tx,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	googleCredentials []byte,
	iapTestOn bool,
) error {
	orders, err := storage.PendingPlaystoreOrders(tx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		machine, err := playstore.New(
			googleCredentials,
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
		if err := machine.Order().UpdateState(tx); err != nil {
			return err
		}
	}

	return nil
}
