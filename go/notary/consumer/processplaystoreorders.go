package consumer

import (
	"crypto/ecdsa"
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/contracts"
)

func ProcessPlaystoreOrders(
	tx *sql.Tx,
	constracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
) error {
	return nil
}
