package consumer

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"

	log "github.com/sirupsen/logrus"
)

func SubmitPlayStorePlayerPurchase(
	tx *sql.Tx,
	in input.SubmitPlayStorePlayerPurchaseInput,
) error {
	log.Debugf("SubmitPlayStorePlayerPurchase %+v", in)

	data, err := playstore.DataFromReceipt(in.Receipt)
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

	service := postgres.NewPlaystoreOrderHistoryService(postgres.NewPlaystoreOrderService(tx))
	if err := service.Insert(order); err != nil {
		return err
	}

	return nil
}
