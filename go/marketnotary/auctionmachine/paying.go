package auctionmachine

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/bidmachine"
	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

func (m *AuctionMachine) processPaying() error {

	if m.Auction.State != storage.AUCTION_PAYING {
		return errors.New("Paying: wrong state")
	}

	idx := bidmachine.IndexFirstAlive(m.Bids)
	if idx == -1 {
		return nil
	}
	bidMachine, err := bidmachine.New(
		m.Auction,
		m.Bids[idx],
		m.market,
		m.freeverse,
		m.client,
		m.db,
	)
	if err != nil {
		return err
	}

	m.Bids[idx], err = bidMachine.Process()
	if err != nil {
		return err
	}
	if m.Bids[idx].State == storage.BID_PAYING {
		return nil
	}
	if m.Bids[idx].State == storage.BID_FAILED_TO_PAY {
		m.Auction.State = storage.AUCTION_FAILED_TO_PAY
	}
	m.Auction.State = storage.AUCTION_PAID

	return nil
}
