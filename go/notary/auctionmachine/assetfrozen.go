package auctionmachine

import (
	"errors"
	"time"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (m *AuctionMachine) ProcessAssetFrozen() error {
	if m.auction.State != storage.AuctionAssetFrozen {
		return errors.New("AssetFrozen: wrong state")
	}

	now := time.Now().Unix()
	if now > m.auction.ValidUntil {
		m.auction.State = storage.AuctionPaying
	}

	return nil
}
