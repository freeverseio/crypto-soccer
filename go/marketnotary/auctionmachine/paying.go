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

		bid, err := bidMachine.Process()
		if err != nil {
			return err
		}
		if bid.State == storage.BID_PAYING {
			return nil
		}
		if bid.State == storage.BID_FAILED_TO_PAY {
			m.Auction.State = storage.AUCTION_FAILED_TO_PAY
		}

		m.Bids[idx] = bid
		m.Auction.State = storage.AUCTION_PAID
	}

	return nil
}
