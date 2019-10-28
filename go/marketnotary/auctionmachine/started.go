package auctionmachine

import (
	"errors"
	"time"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

type Started struct {
}

func NewStarted() State {
	return &Started{}
}

func (b *Started) Process(m *AuctionMachine) error {
	if m.Auction.State != storage.AUCTION_STARTED {
		return errors.New("Started: wrong state")
	}
	now := time.Now().Unix()

	if (len(m.Bids) == 0) && (m.Auction.ValidUntil.Int64()) < now {
		m.Auction.State = storage.AUCTION_NO_BIDS
		m.SetState(NewNoBids())
	}

	return nil
}
