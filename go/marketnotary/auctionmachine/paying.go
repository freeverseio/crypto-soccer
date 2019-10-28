package auctionmachine

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

type Paying struct {
}

func NewPaying() State {
	return &Paying{}
}

func (b *Paying) Process(m *AuctionMachine) error {
	if m.Auction.State != storage.AUCTION_PAYING {
		return errors.New("Paying: wrong state")
	}

	return nil
}
