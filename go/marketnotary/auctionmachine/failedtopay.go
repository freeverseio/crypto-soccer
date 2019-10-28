package auctionmachine

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

type FailedToPay struct {
}

func NewFailedToPay() State {
	return &FailedToPay{}
}

func (b *FailedToPay) Process(m *AuctionMachine) error {
	if m.Auction.State != storage.AUCTION_FAILED_TO_PAY {
		return errors.New("FailedToPay: wrong state")
	}

	return nil
}
