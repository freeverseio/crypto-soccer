package auctionmachine

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	log "github.com/sirupsen/logrus"
)

type AuctionMachine struct {
	auction   storage.Auction
	bids      []storage.Bid
	contracts contracts.Contracts
	freeverse *ecdsa.PrivateKey
}

func New(
	auction storage.Auction,
	bids []storage.Bid,
	contracts contracts.Contracts,
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
	}, nil
}

func (b *AuctionMachine) Process(market marketpay.IMarketPay) error {
	log.Infof("Process auction %v in state %v", b.auction.ID, b.State())
	switch b.auction.State {
	case storage.AuctionStarted:
		return b.processStarted()
	case storage.AuctionAssetFrozen:
		return b.ProcessAssetFrozen()
	case storage.AuctionPaying:
		return b.ProcessPaying(market)
	case storage.AuctionWithdrableBySeller:
		log.Warn("auctionmachine AuctionWithdrabeBySeller not implemented")
		return nil
	case storage.AuctionWithdrableByBuyer:
		log.Warn("auctionmachine AuctionWithdrabeByBuyer not implemented")
		return nil
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
	return b.auction.State
}

func (b AuctionMachine) StateExtra() string {
	return b.auction.StateExtra
}

func (b AuctionMachine) Auction() storage.Auction {
	return b.auction
}

func (b AuctionMachine) Bids() []storage.Bid {
	return b.bids
}
