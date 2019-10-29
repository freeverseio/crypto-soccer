package auctionmachine

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
	log "github.com/sirupsen/logrus"
)

type FailedToFreeze struct {
}

func NewFailedToFreeze() State {
	return &FailedToFreeze{}
}

func (b *FailedToFreeze) Process(m *AuctionMachine) error {
	if m.Auction.State != storage.AUCTION_FAILED_TO_FREEZE {
		return errors.New("FailedToFreeze: wrong state")
	}

	log.Warn("FailedToFreeze::Process called")
	return nil
}
