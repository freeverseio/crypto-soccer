package auctionmachine

import (
	"errors"
	"time"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
	log "github.com/sirupsen/logrus"
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

	if len(m.Bids) == 0 {
		if now > m.Auction.ValidUntil.Int64() {
			m.Auction.State = storage.AUCTION_NO_BIDS
			m.SetState(NewNoBids())
		}
		return nil
	}

	// TODO trying to freeze the asset
	log.Warn("Started::Process| Not implemented Auctrion going ASSET_FROZEN")
	m.Auction.State = storage.AUCTION_ASSET_FROZEN
	m.SetState(NewAssetFrozen())

	return nil
}
