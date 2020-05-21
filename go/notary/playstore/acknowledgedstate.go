package playstore

import (
	"context"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	log "github.com/sirupsen/logrus"
)

func (b *Machine) processAcknowledged(ctx context.Context) error {
	purchase, err := b.client.VerifyProduct(
		ctx,
		b.order.PackageName,
		b.order.ProductId,
		b.order.PurchaseToken,
	)
	if err != nil {
		b.setState(storage.PlaystoreOrderFailed, err.Error())
		return nil
	}

	validator := NewPurchaseValidator(*purchase)
	if validator.IsTest() && !b.iapTestOn {
		log.Warningf("[consumer|iap] received test orderId %v ... skip creating player", purchase.OrderId)
	} else {
		if err := b.assignAsset(); err != nil {
			b.setState(storage.PlaystoreOrderRefunding, err.Error())
			return nil
		}
		log.Infof("[consumer|iap] orderId %v playerId %v assigned to teamId %v", purchase.OrderId, b.order.PlayerId, b.order.TeamId)
	}

	b.setState(storage.PlaystoreOrderComplete, "")
	return nil
}
