package auctionmachine

import (
	"errors"
	"time"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (m *AuctionMachine) processAssetFrozen() error {
	if m.Auction.State != storage.AUCTION_ASSET_FROZEN {
		return errors.New("AssetFrozen: wrong state")
	}

	now := time.Now().Unix()
	if now > m.Auction.ValidUntil.Int64() {
		m.Bids[0].State = storage. BIDPAYING
		m.Auction.State = storage.AUCTION_PAYING
	}

	return nil
}
