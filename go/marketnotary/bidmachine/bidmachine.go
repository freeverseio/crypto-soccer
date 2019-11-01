package bidmachine

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

type BidMachine struct {
	auction storage.Auction
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

func (b *BidMachine) HasAliveBids() bool {
	for _, bid := range b.Bids {
		if (bid.State == storage.BID_ACCEPTED) ||
			(bid.State == storage.BID_PAYING) {
			return true
		}
	}
	return false
}

func (b *BidMachine) IndexFirstAlive() int {
	if b.HasAliveBids() == false {
		return -1
	}

	// var idx = 0
	// for i, bid := range b.Bids {
	// 	if bid.State == storage.BID_ACCEPTED {
	// 		return idx
	// 	}
	// }
	return -1
}

func (b *BidMachine) Process() error {
	idx := b.IndexFirstAlive()
	if idx == -1 {
		return nil
	}

	return nil
}
