package auctionmachine

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/marketpay"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	log "github.com/sirupsen/logrus"
)

type AuctionMachine struct {
	auction              storage.Auction
	bids                 []storage.Bid
	contracts            contracts.Contracts
	freeverse            *ecdsa.PrivateKey
	shouldQueryMarketPay bool
}

func New(
	auction storage.Auction,
	bids []storage.Bid,
	contracts contracts.Contracts,
	freeverse *ecdsa.PrivateKey,
	shouldQueryMarketPay bool,
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
		shouldQueryMarketPay,
	}, nil
}

func (b *AuctionMachine) Process(market marketpay.MarketPayService) error {
	log.Debugf("Process auction %v in state %v", b.auction.ID, b.State())
	switch b.auction.State {
	case storage.AuctionStarted:
		return b.processStarted()
	case storage.AuctionAssetFrozen:
		return b.ProcessAssetFrozen()
	case storage.AuctionPaying:
		return b.ProcessPaying(market)
	case storage.AuctionWithdrableBySeller:
		return b.ProcessWithdrawableBySeller(market)
	case storage.AuctionWithdrableByBuyer:
		log.Warn("auctionmachine AuctionWithdrabeByBuyer not implemented")
		return nil
	case storage.AuctionValidation:
		return b.ProcessValidation(market)
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

func (b *AuctionMachine) SetState(state storage.AuctionState, extra string) {
	if state == storage.AuctionFailed {
		log.Warnf("auction %v in state %v with %v", b.auction.ID, state, extra)
	}
	b.auction.State = state
	b.auction.StateExtra = extra
}

func (b AuctionMachine) checkState(state storage.AuctionState) error {
	if b.auction.State != state {
		return fmt.Errorf("auction[%v|%v] is not in state %v", b.auction.ID, b.auction.State, state)
	}
	return nil
}
