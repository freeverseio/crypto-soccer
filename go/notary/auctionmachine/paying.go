package auctionmachine

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/bidmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func (m *AuctionMachine) processPaying() error {

	if m.Auction.State != storage.AUCTION_PAYING {
		return errors.New("Paying: wrong state : " + string(m.Auction.State))
	}

	bid := bidmachine.FirstAlive(m.Bids)
	if bid == nil {
		m.Auction.State = storage.AUCTION_FAILED
		m.Auction.StateExtra = "Failed to pay"
		return nil
	}

	bidMachine, err := bidmachine.New(
		m.Auction,
		bid,
		m.market,
		m.freeverse,
		m.client,
	)
	if err != nil {
		return err
	}

	err = bidMachine.Process()
	if err != nil {
		return err
	}
	if bid.State == storage.BIDPAID {
		log.Infof("[auction] %v PAYING -> PAID", m.Auction.UUID)
		m.Auction.State = storage.AUCTION_PAID
	}

	return nil
}
