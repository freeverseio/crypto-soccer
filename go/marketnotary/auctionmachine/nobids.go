package auctionmachine

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
	log "github.com/sirupsen/logrus"
)

func (m *AuctionMachine) ProcessNoBids() error {
	if m.Auction.State != storage.AUCTION_NO_BIDS {
		return errors.New("NoBids: wrong state")
	}

	log.Warn("NoBids::Process called")
	return nil
}
