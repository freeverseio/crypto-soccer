package auctionmachine

import (
	"errors"
	"time"

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

	now := time.Now().Unix()
	if (now - m.Auction.ValidUntil.Int64()) > 60 {
		if now%2 == 0 {
			m.Auction.State = storage.AUCTION_PAID
			m.SetState(NewPaid())
		} else {
			m.Auction.State = storage.AUCTION_FAILED_TO_PAY
			m.SetState(NewFailedToPay())
		}
	}

	return nil
}
