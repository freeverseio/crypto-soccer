package playstore

import (
	"context"
	"fmt"
	"time"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/worldplayer"
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

	worldPlayerService := worldplayer.NewWorldPlayerService(b.contracts, b.namesdb)
	worldPlayer, err := worldPlayerService.GetWorldPlayer(
		string(b.order.PlayerId),
		string(b.order.TeamId),
		time.Now().Unix(),
	)
	if err != nil {
		b.setState(storage.PlaystoreOrderOpen, err.Error())
		return nil
	}
	if worldPlayer == nil {
		b.setState(storage.PlaystoreOrderRefunding, fmt.Sprintf("orderId %v has an invalid playerId %v", b.order.OrderId, b.order.PlayerId))
		return nil
	}
	if worldPlayer.ProductId() != b.order.ProductId {
		b.setState(storage.PlaystoreOrderRefunding, fmt.Sprintf("orderId %v has an productId mismatch %v != %v", b.order.OrderId, worldPlayer.ProductId(), b.order.ProductId))
		return nil
	}

	payload := fmt.Sprintf("playerId: %v", b.order.PlayerId)
	if err := b.client.AcknowledgePurchase(
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
