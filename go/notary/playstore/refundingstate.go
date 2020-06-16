package playstore

import (
	"context"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *Machine) processRefundingState(ctx context.Context) error {
	err := b.client.Refund(ctx, b.order.PackageName, b.order.OrderId)
	if err != nil {
		b.setState(storage.PlaystoreOrderRefunding, err.Error())
	}
	b.setState(storage.PlaystoreOrderRefunded, "")
	return nil
}
