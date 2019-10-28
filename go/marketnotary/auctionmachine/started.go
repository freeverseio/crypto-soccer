package auctionmachine

import (
	"time"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

type Started struct {
}

func NewStarted() State {
	return &Started{}
}

func (b *Started) Process(m *AuctionMachine) {
	now := time.Now().Unix()

	if (len(m.Bids) == 0) && (m.Auction.ValidUntil.Int64()) < now {
		m.Auction.State = storage.AUCTION_NO_BIDS
	}

}
