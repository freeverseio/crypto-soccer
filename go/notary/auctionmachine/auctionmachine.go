package auctionmachine

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

type AuctionMachine struct {
	Auction   storage.Auction
	Bids      []*storage.Bid
	contracts *contracts.Contracts
	freeverse *ecdsa.PrivateKey
	signer    *signer.Signer
}

func New(
	auction storage.Auction,
	bids []*storage.Bid,
	contracts *contracts.Contracts,
	freeverse *ecdsa.PrivateKey,
) (*AuctionMachine, error) {
	if contracts.Market == nil {
		return nil, errors.New("market is nil")
	}
	if freeverse == nil {
		return nil, errors.New("owner is nil")
	}
	return &AuctionMachine{
		auction,
		bids,
		contracts,
		freeverse,
		signer.NewSigner(contracts, freeverse),
	}, nil
}

func (b *AuctionMachine) Process(market marketpay.IMarketPay) error {
	switch b.Auction.State {
	case storage.AUCTION_STARTED:
		return b.processStarted()
	case storage.AUCTION_ASSET_FROZEN:
		return b.processAssetFrozen()
	case storage.AUCTION_PAYING:
		return b.processPaying(market)
	case storage.AuctionCancelled:
		return nil
	case storage.AuctionFailed:
		return nil
	case storage.AuctionEnded:
		return nil
	default:
		return fmt.Errorf("Unknown auction state %v", b.State())
	}
}

func (b AuctionMachine) State() storage.AuctionState {
	return b.Auction.State
}
