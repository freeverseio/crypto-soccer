package bidmachine

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

type BidMachine struct {
	Auction storage.Auction
	Bids    []storage.Bid
}

func New(
	auction storage.Auction,
	bids []storage.Bid,
) (*BidMachine, error) {
	if auction.State != storage.AUCTION_PAYING {
		return nil, errors.New("Auction is not in PAYING state")
	}
	return &BidMachine{
		auction,
		bids,
	}, nil
}
