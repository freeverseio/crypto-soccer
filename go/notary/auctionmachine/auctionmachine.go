package auctionmachine

import (
	"crypto/ecdsa"
	"errors"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	log "github.com/sirupsen/logrus"
)

type AuctionMachine struct {
	Auction   *storage.Auction
	Bids      []*storage.Bid
	contracts *contracts.Contracts
	freeverse *ecdsa.PrivateKey
	signer    *signer.Signer
}

func New(
	auction *storage.Auction,
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

func (b *AuctionMachine) Process() error {
	switch b.Auction.State {
	case storage.AUCTION_STARTED:
		return b.processStarted()
	case storage.AUCTION_ASSET_FROZEN:
		return b.processAssetFrozen()
	case storage.AUCTION_PAYING:
		return b.processPaying()
	default:
		return b.processUnknownState()
	}
}

func (b *AuctionMachine) processUnknownState() error {
	log.Infof("[auction] %v: unknown state %v", b.Auction.UUID, b.Auction.State)
	b.Auction.StateExtra = "Unknown state " + string(b.Auction.State)
	b.Auction.State = storage.AUCTION_FAILED
	return nil
}
