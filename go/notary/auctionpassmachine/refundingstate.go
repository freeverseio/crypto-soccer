package auctionpassmachine

import (
	"context"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *AuctionPassMachine) processAuctionPassRefundingState(ctx context.Context) error {
	err := b.client.Refund(ctx, b.order.PackageName, b.order.OrderId)
	if err != nil {
		b.setState(storage.AuctionPassPlaystoreOrderRefunding, err.Error())
		return nil
	}
	b.setState(storage.AuctionPassPlaystoreOrderRefunded, "")
	return nil
}
