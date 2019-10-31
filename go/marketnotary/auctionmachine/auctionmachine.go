package auctionmachine

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/ethclient"
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
	market    *market.Market
	freeverse *ecdsa.PrivateKey
	signer    *signer.Signer
	client    *ethclient.Client
}

func New(
	auction storage.Auction,
	bids []storage.Bid,
	market *market.Market,
	freeverse *ecdsa.PrivateKey,
	client *ethclient.Client,
) (*AuctionMachine, error) {
	if market == nil {
		return nil, errors.New("market is nil")
	}
	if freeverse == nil {
		return nil, errors.New("owner is nil")
	}
	if client == nil {
		return nil, errors.New("client is nil")
	}
	return &AuctionMachine{
		auction,
		bids,
		market,
		freeverse,
		signer.NewSigner(market, freeverse),
		client,
	}, nil
}

func (b *AuctionMachine) Process() error {
	var state State
	switch b.Auction.State {
	case storage.AUCTION_STARTED:
		state = NewStarted()
	case storage.AUCTION_ASSET_FROZEN:
		state = NewAssetFrozen()
	case storage.AUCTION_PAYING:
		state = NewPaying()
	case storage.AUCTION_NO_BIDS:
		state = NewNoBids()
	default:
		return errors.New("unknown auction state")
	}
	return state.Process(b)
}
