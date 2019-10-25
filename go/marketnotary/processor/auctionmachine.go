package processor

import (
	"time"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

type AuctionMachine struct {
	Auction storage.Auction
	current State
}

func NewAuctionMachine(auction storage.Auction) *AuctionMachine {
	var state State
	switch auction.State {
	case storage.AUCTION_STARTED:
		state = NewStarted()
	}
	return &AuctionMachine{auction, state}
}

func (b *AuctionMachine) SetCurrent(s State) {
	b.current = s
}

func (b *AuctionMachine) Process() {
	b.current.Process(b)
}

type State interface {
	Process(m *AuctionMachine)
}

type Started struct {
}

func NewStarted() State {
	return &Started{}
}

func (b *Started) Process(m *AuctionMachine) {
	now := time.Now().Unix()
	if m.Auction.ValidUntil.Int64() < now {
		m.Auction.State = storage.AUCTION_NO_BIDS
	}
}
