package auctionmachine

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

type Paid struct {
}

func NewPaid() State {
	return &Paid{}
}

func (b *Paid) Process(m *AuctionMachine) error {
	if m.Auction.State != storage.AUCTION_PAID {
		return errors.New("Paid: wrong state")
	}

	return nil
}
