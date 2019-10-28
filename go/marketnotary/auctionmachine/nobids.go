package auctionmachine

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

type NoBids struct {
}

func NewNoBids() State {
	return &NoBids{}
}

func (b *NoBids) Process(m *AuctionMachine) error {
	if m.Auction.State != storage.AUCTION_NO_BIDS {
		return errors.New("NoBids: wrong state")
	}

	return nil
}
