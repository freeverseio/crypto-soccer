package auctionmachine

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go/contracts/market"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	log "github.com/sirupsen/logrus"
)

type AuctionMachine struct {
	Auction   *storage.Auction
	Bids      []*storage.Bid
	market    *market.Market
	freeverse *ecdsa.PrivateKey
	signer    *signer.Signer
	client    *ethclient.Client
}

func New(
	auction *storage.Auction,
	bids []*storage.Bid,
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
	switch b.Auction.State {
	case storage.AUCTION_STARTED:
		return b.processStarted()
	case storage.AUCTION_ASSET_FROZEN:
		return b.processAssetFrozen()
	case storage.AUCTION_PAYING:
		return b.processPaying()
	case storage.AUCTION_NO_BIDS:
		return b.processNoBids()
	default:
		return b.processUnknownState()
	}
}

func (b *AuctionMachine) processUnknownState() error {
	log.Infof("[auction] %v: unknown state %v", b.Auction.UUID, b.Auction.State)

	return nil
}
