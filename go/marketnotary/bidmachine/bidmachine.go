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

func IndexFirstAlive(bids []storage.Bid) int {
	// first searching for PAYING bid
	for i, bid := range bids {
		if bid.State == storage.BID_PAYING {
			return i
		}
	}
	// then search for the highest ACCEPTED bid
	idx := -1
	extraPrice := int64(-1)
	for i, bid := range bids {
		if bid.State == storage.BID_ACCEPTED {
			if idx == -1 {
				idx = i
				extraPrice = bid.ExtraPrice
			} else {
				if bid.ExtraPrice > extraPrice {
					idx = i
					extraPrice = bid.ExtraPrice
				}
			}
		}
	}
	return idx
}

func (b *BidMachine) Process() error {
	idx := IndexFirstAlive(b.Bids)
	if idx == -1 {
		return nil
	}
	bid := b.Bids[idx]
	switch bid.State {
	case storage.BID_PAYING:
		return b.processPaying(idx)
	case storage.BID_ACCEPTED:
		return b.processAccepted(idx)
	default:
		return errors.New("Unknown bid state")
	}
}

func (b *BidMachine) processPaying(idx int) error {
	b.Bids[idx].State = storage.BID_PAID
	return nil
}

func (b *BidMachine) processAccepted(idx int) error {
	b.Bids[idx].State = storage.BID_PAYING
	return nil
}
