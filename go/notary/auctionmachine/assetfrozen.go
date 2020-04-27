package auctionmachine

import (
	"time"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *AuctionMachine) ProcessAssetFrozen() error {
	if err := b.checkState(storage.AuctionAssetFrozen); err != nil {
		return err
	}

	now := time.Now().Unix()
	if now > b.auction.ValidUntil {
		b.auction.State = storage.AuctionPaying
	}

	return nil
}
