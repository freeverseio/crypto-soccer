package bidmachine

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go/contracts/market"
	"github.com/freeverseio/crypto-soccer/go/helper"
	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	log "github.com/sirupsen/logrus"
)

type BidMachine struct {
	auction   *storage.Auction
	bid       *storage.Bid
	market    *market.Market
	freeverse *ecdsa.PrivateKey
	signer    *signer.Signer
	client    *ethclient.Client
}

func New(
	auction *storage.Auction,
	bid *storage.Bid,
	market *market.Market,
	freeverse *ecdsa.PrivateKey,
	client *ethclient.Client,
) (*BidMachine, error) {
	if auction.State != storage.AUCTION_PAYING {
		return nil, errors.New("Auction is not in PAYING state")
	}
	if auction.UUID != bid.Auction {
		return nil, errors.New("Bid of wrong auction")
	}
	return &BidMachine{
		auction,
		bid,
		market,
		freeverse,
		signer.NewSigner(market, freeverse),
		client,
	}, nil
}

func IndexFirstAlive(bids []*storage.Bid) int {
	// first searching for PAYING bid
	for i, bid := range bids {
		if bid.State == storage.BIDPAYING {
			return i
		}
	}
	// then search for the highest ACCEPTED bid
	idx := -1
	extraPrice := int64(-1)
	for i, bid := range bids {
		if bid.State == storage.BIDACCEPTED {
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
	switch b.bid.State {
	case storage.BIDPAYING:
		return b.processPaying()
	case storage.BIDACCEPTED:
		return b.processAccepted()
	case storage.BIDFAILEDTOPAY:
		return nil
	default:
		return errors.New("Unknown bid state")
	}
}

func (b *BidMachine) processPaying() error {
	if b.bid.PaymentID == "" {
		log.Infof("[bid] Auction %v, extra_price %v | create MarketPay order", b.bid.Auction, b.bid.ExtraPrice)
		market, err := marketpay.New()
		if err != nil {
			return err
		}
		price := fmt.Sprintf("%.2f", float64(b.auction.Price.Int64()+b.bid.ExtraPrice)/100.0)
		name := "Freeverse Player transaction"
		order, err := market.CreateOrder(name, price)
		if err != nil {
			return err
		}
		b.bid.PaymentID = order.TrusteeShortlink.Hash
		b.bid.PaymentURL = order.TrusteeShortlink.ShortURL
	} else {
		log.Warningf("[bid] Auction %v, extra_price %v | waiting for order %v to be processed", b.bid.Auction, b.bid.ExtraPrice, b.bid.PaymentID)
		market, err := marketpay.New()
		if err != nil {
			return err
		}
		order, err := market.GetOrder(b.bid.PaymentID)
		if err != nil {
			return err
		}
		paid := market.IsPaid(*order)
		if paid {
			isOffer2StartAuction := false
			bidHiddenPrice, err := b.signer.BidHiddenPrice(big.NewInt(b.bid.ExtraPrice), big.NewInt(b.bid.Rnd))
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
				big.NewInt(b.bid.ExtraPrice),
				big.NewInt(b.bid.Rnd),
				b.bid.TeamID,
				isOffer2StartAuction,
			)
			if err != nil {
				return err
			}
			sig[1], sig[2], sigV, err = b.signer.RSV(b.bid.Signature)
			if err != nil {
				return err
			}
			tx, err := b.market.CompletePlayerAuction(
				bind.NewKeyedTransactor(b.freeverse),
				auctionHiddenPrice,
				b.auction.ValidUntil,
				b.auction.PlayerID,
				bidHiddenPrice,
				b.bid.TeamID,
				sig,
				sigV,
				isOffer2StartAuction,
			)
			if err != nil {
				b.bid.State = storage.BIDFAILED
				b.bid.StateExtra = err.Error()
				return err
			}
			receipt, err := helper.WaitReceipt(b.client, tx, 60)
			if err != nil {
				b.bid.State = storage.BIDFAILED
				b.bid.StateExtra = "Timeout waiting for the receipt"
				return err
			}
			if receipt.Status == 0 {
				b.bid.State = storage.BIDFAILED
				b.bid.StateExtra = "Mined but receipt.Status == 0"
				return err
			}
			b.auction.PaymentURL = order.SettlorShortlink.ShortURL
			b.bid.State = storage.BIDPAID
		}
	}
	return nil
}

func (b *BidMachine) processAccepted() error {
	b.bid.State = storage.BIDPAYING
	return nil
}
