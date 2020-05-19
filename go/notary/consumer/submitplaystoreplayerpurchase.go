package consumer

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"google.golang.org/api/androidpublisher/v3"

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

func isTestPurchase(purchase *androidpublisher.ProductPurchase) bool {
	if purchase.PurchaseType != nil {
		if *purchase.PurchaseType == 0 { // Test
			return true
		}
	}
	return false
}
