package auctionpassmachine

import (
	"context"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/notary/googleplaystoreutils"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *AuctionPassMachine) processAuctionPassOpenState(ctx context.Context, service storage.Tx) error {
	purchase, err := b.client.GetPurchase(
		ctx,
		b.order.PackageName,
		b.order.ProductId,
		b.order.PurchaseToken,
	)
	if err != nil {
		b.setState(storage.AuctionPassPlaystoreOrderFailed, err.Error())
		return nil
	}

	validator := googleplaystoreutils.NewPurchaseValidator(*purchase)
	if validator.IsCanceled() {
		b.setState(storage.AuctionPassPlaystoreOrderComplete, "cancelled")
		return nil
	}
	if validator.IsPending() {
		b.setState(storage.AuctionPassPlaystoreOrderOpen, "pending")
		return nil
	}
	if !validator.IsPurchased() {
		b.setState(storage.AuctionPassPlaystoreOrderFailed, "invalid puchase state")
		return nil
	}
	if validator.IsAcknowledged() {
		b.setState(storage.AuctionPassPlaystoreOrderFailed, "already acknowledged")
		return nil
	}

	payload := fmt.Sprintf("owner: %v", b.order.Owner)
	if err := b.client.AcknowledgePurchase(
		ctx,
		b.order.PackageName,
		b.order.ProductId,
		b.order.PurchaseToken,
		payload,
	); err != nil {
		b.setState(storage.AuctionPassPlaystoreOrderOpen, err.Error())
		return nil
	}

	b.setState(storage.AuctionPassPlaystoreOrderAcknowledged, "")
	return nil
}
