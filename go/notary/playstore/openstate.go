package playstore

import (
	"context"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *Machine) processOpenState(ctx context.Context) error {
	purchase, err := b.client.GetPurchase(
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
	if validator.IsCanceled() {
		b.setState(storage.PlaystoreOrderComplete, "cancelled")
		return nil
	}
	if validator.IsPending() {
		b.setState(storage.PlaystoreOrderOpen, "pending")
		return nil
	}
	if !validator.IsPurchased() {
		b.setState(storage.PlaystoreOrderFailed, "invalid puchase state")
		return nil
	}
	if validator.IsAcknowledged() {
		b.setState(storage.PlaystoreOrderFailed, "already acknowledged")
		return nil
	}

	payload := fmt.Sprintf("playerId: %v", b.order.PlayerId)
	if err := b.client.AcknowledgedPurchase(
		ctx,
		b.order.PackageName,
		b.order.ProductId,
		b.order.PurchaseToken,
		payload,
	); err != nil {
		b.setState(storage.PlaystoreOrderOpen, err.Error())
		return nil
	}

	b.setState(storage.PlaystoreOrderAcknowledged, "")
	return nil
}
