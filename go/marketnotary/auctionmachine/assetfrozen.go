package auctionmachine

import (
	"errors"
	"time"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

type AssetFrozen struct {
}

func NewAssetFrozen() State {
	return &AssetFrozen{}
}

func (b *AssetFrozen) Process(m *AuctionMachine) error {
	if m.Auction.State != storage.AUCTION_ASSET_FROZEN {
		return errors.New("AssetFrozen: wrong state")
	}

	now := time.Now().Unix()
	if now > m.Auction.ValidUntil.Int64() {
		m.Auction.State = storage.AUCTION_PAYING
	}

	return nil
}
