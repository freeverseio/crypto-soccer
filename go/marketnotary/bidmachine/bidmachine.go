package bidmachine

import (
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go/contracts/market"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/marketnotary/signer"
	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

type BidMachine struct {
	auction   storage.Auction
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
	signer *signer.Signer,
	client *ethclient.Client,
) (*BidMachine, error) {
	if auction.State != storage.AUCTION_PAYING {
		return nil, errors.New("Auction is not in PAYING state")
	}
	return &BidMachine{
		auction,
		bids,
		market,
		freeverse,
		signer,
		client,
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
	bid := b.Bids[idx]
	isOffer2StartAuction := false
	bidHiddenPrice, err := b.signer.BidHiddenPrice(big.NewInt(bid.ExtraPrice), big.NewInt(bid.Rnd))
	if err != nil {
		return err
	}
	auctionHiddenPrice, err := b.signer.HashPrivateMsg(b.auction.CurrencyID, b.auction.Price, b.auction.Rnd)
	if err != nil {
		return err
	}
	var sig [3][32]byte
	var sigV uint8
	sig[0], err = b.signer.HashBidMessage(
		b.auction.CurrencyID,
		b.auction.Price,
		b.auction.Rnd,
		b.auction.ValidUntil,
		b.auction.PlayerID,
		big.NewInt(bid.ExtraPrice),
		big.NewInt(bid.Rnd),
		bid.TeamID,
		isOffer2StartAuction,
	)
	if err != nil {
		return err
	}
	sig[1], sig[2], sigV, err = b.signer.RSV(bid.Signature)
	if err != nil {
		return err
	}
	tx, err := b.market.CompletePlayerAuction(
		bind.NewKeyedTransactor(b.freeverse),
		auctionHiddenPrice,
		b.auction.ValidUntil,
		b.auction.PlayerID,
		bidHiddenPrice,
		bid.TeamID,
		sig,
		sigV,
		isOffer2StartAuction,
	)
	if err != nil {
		b.Bids[idx].State = storage.BID_FAILED_TO_PAY
		return err
	}
	receipt, err := helper.WaitReceipt(b.client, tx, 60)
	if err != nil {
		b.Bids[idx].State = storage.BID_FAILED_TO_PAY
		return err
	}
	if receipt.Status == 0 {
		b.Bids[idx].State = storage.BID_FAILED_TO_PAY
		return err
	}

	b.Bids[idx].State = storage.BID_PAID
	return nil
}

func (b *BidMachine) processAccepted(idx int) error {
	b.Bids[idx].State = storage.BID_PAYING
	return nil
}
