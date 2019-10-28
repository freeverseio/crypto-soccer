package auctionmachine

import (
	"github.com/freeverseio/crypto-soccer/go/contracts/market"
	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

type State interface {
	Process(m *AuctionMachine) error
}

type AuctionMachine struct {
	Auction storage.Auction
	Bids    []storage.Bid
	current State
	market  *market.Market
}

func NewAuctionMachine(
	auction storage.Auction,
	bids []storage.Bid,
	market *market.Market,
) *AuctionMachine {
	var state State
	switch auction.State {
	case storage.AUCTION_STARTED:
		state = NewStarted()
	}
	return &AuctionMachine{
		auction,
		bids,
		state,
		market,
	}
}

func (b *AuctionMachine) Process() error {
	return b.current.Process(b)
}

func (b *AuctionMachine) SetState(state State) error {
	b.current = state
	return nil
}
