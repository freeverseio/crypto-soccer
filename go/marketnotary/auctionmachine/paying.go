package auctionmachine

import (
	"errors"
	"time"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/bidmachine"
	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

func (m *AuctionMachine) processPaying() error {
	if m.Auction.State != storage.AUCTION_PAYING {
		return errors.New("Paying: wrong state")
	}

	now := time.Now().Unix()
	if (now - m.Auction.ValidUntil.Int64()) > 2 {
		bidMachine, err := bidmachine.New(
			m.Auction,
			m.Bids,
			m.market,
			m.freeverse,
			m.signer,
			m.client,
		)
		if err != nil {
			return err
		}

		err = bidMachine.Process()
		if err != nil {
			m.Auction.State = storage.AUCTION_FAILED_TO_PAY
			return err
		}

		m.Auction.State = storage.AUCTION_PAID
	}

	return nil
}
