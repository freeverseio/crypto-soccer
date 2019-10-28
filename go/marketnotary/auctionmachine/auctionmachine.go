package auctionmachine

import (
	"crypto/ecdsa"
	"errors"

	"github.com/freeverseio/crypto-soccer/go/contracts/market"
	"github.com/freeverseio/crypto-soccer/go/marketnotary/signer"
	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

type State interface {
	Process(m *AuctionMachine) error
}

type AuctionMachine struct {
	Auction   storage.Auction
	Bids      []storage.Bid
	current   State
	market    *market.Market
	freeverse *ecdsa.PrivateKey
	signer    *signer.Signer
}

func New(
	auction storage.Auction,
	bids []storage.Bid,
	market *market.Market,
	freeverse *ecdsa.PrivateKey,
) (*AuctionMachine, error) {
	var state State
	switch auction.State {
	case storage.AUCTION_STARTED:
		state = NewStarted()
	case storage.AUCTION_ASSET_FROZEN:
		state = NewAssetFrozen()
	case storage.AUCTION_PAYING:
		state = NewPaying()
	default:
		return nil, errors.New("unknown auction state")
	}
	return &AuctionMachine{
		auction,
		bids,
		state,
		market,
		freeverse,
		signer.NewSigner(market),
	}, nil
}

func (b *AuctionMachine) Process() error {
	return b.current.Process(b)
}

func (b *AuctionMachine) SetState(state State) error {
	b.current = state
	return nil
}
