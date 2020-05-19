package consumer

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	log "github.com/sirupsen/logrus"
)

func SubmitPlayStorePlayerPurchase(
	tx *sql.Tx,
	in input.SubmitPlayStorePlayerPurchaseInput,
) error {
	log.Debugf("SubmitPlayStorePlayerPurchase %+v", in)

	data, err := playstore.InappPurchaseDataFromReceipt(in.Receipt)
	if err != nil {
		return err
	}

	order := storage.NewPlaystoreOrder()
	order.OrderId = data.OrderId
	order.PackageName = data.PackageName
	order.ProductId = data.ProductId
	order.PurchaseToken = data.PurchaseToken
	order.PlayerId = string(in.PlayerId)
	order.TeamId = string(in.TeamId)
	order.Signature = in.Signature
	if err := order.Insert(tx); err != nil {
		return err
	}

	return nil
}
