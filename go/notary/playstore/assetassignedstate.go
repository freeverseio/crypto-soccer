package playstore

import (
	"context"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *Machine) processAssetAssigned(ctx context.Context) error {
	payload := fmt.Sprintf("playerId: %v", b.order.PlayerId)
	if err := b.client.AcknowledgeProduct(
		ctx,
		b.order.PackageName,
		b.order.ProductId,
		b.order.PurchaseToken,
		payload,
	); err != nil {
		b.setState(storage.PlaystoreOrderAssetAssigned, err.Error())
		return err
	}

	b.setState(storage.PlaystoreOrderComplete, "")
	return nil
}
