package auctionmachine

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/bidmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
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
	)
	if err != nil {
		return err
	}

	err = bidMachine.Process()
	if err != nil {
		return err
	}
	if m.Bids[idx].State == storage.BIDPAYING {
		return nil
	}
	if m.Bids[idx].State == storage.BIDFAILEDTOPAY {
		m.Auction.State = storage.AUCTION_FAILED
		m.Auction.StateExtra = "Faild to pay"
	}
	m.Auction.State = storage.AUCTION_PAID

	return nil
}
