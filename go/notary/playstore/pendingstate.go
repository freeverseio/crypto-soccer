package playstore

import (
	"context"
	"math/big"

	"github.com/awa/go-iap/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func (b *Machine) processPendingState() error {
	playerId, _ := new(big.Int).SetString(b.order.PlayerId, 10)
	if playerId == nil {
		b.setState(storage.PlaystoreOrderFailed, "invalid player")
		return nil
	}
	teamId, _ := new(big.Int).SetString(b.order.TeamId, 10)
	if teamId == nil {
		b.setState(storage.PlaystoreOrderFailed, "invalid team")
		return nil
	}

	client, err := playstore.New(b.googleCredentials)
	if err != nil {
		return err
	}
	ctx := context.Background()

	purchase, err := client.VerifyProduct(
		ctx,
		b.order.PackageName,
		b.order.ProductId,
		b.order.PurchaseToken,
	)
	if err != nil {
		b.setState(storage.PlaystoreOrderFailed, err.Error())
		return nil
	}

	log.Infof("Order %v, state %v, consumed %v, ack %v", purchase.OrderId, purchase.PurchaseState, purchase.ConsumptionState, purchase.AcknowledgementState)

	b.setState(storage.PlaystoreOrderVerified, "")
	return nil
}
