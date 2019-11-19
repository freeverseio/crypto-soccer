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
	"github.com/freeverseio/crypto-soccer/go/marketnotary/signer"
	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"

	log "github.com/sirupsen/logrus"
)

type BidMachine struct {
	auction   storage.Auction
	bid       storage.Bid
	market    *market.Market
	freeverse *ecdsa.PrivateKey
	signer    *signer.Signer
	client    *ethclient.Client
	db        *storage.Storage
}

func New(
	auction storage.Auction,
	bid storage.Bid,
	market *market.Market,
	freeverse *ecdsa.PrivateKey,
	client *ethclient.Client,
	db *storage.Storage,
) (*BidMachine, error) {
	if auction.State != storage.AUCTION_PAYING {
		return nil, errors.New("Auction is not in PAYING state")
	}
	return &BidMachine{
		auction,
		bid,
		market,
		freeverse,
		signer.NewSigner(market, freeverse),
		client,
		db,
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

func (b *BidMachine) Process() (storage.Bid, error) {
	switch b.bid.State {
	case storage.BID_PAYING:
		return b.bid, b.processPaying()
	case storage.BID_ACCEPTED:
		return b.bid, b.processAccepted()
	case storage.BID_FAILED_TO_PAY:
		return b.bid, nil
	default:
		return b.bid, errors.New("Unknown bid state")
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
		err = b.db.UpdateBidPaymentID(b.bid.Auction, b.bid.ExtraPrice, order.TrusteeShortlink.Hash)
		if err != nil {
			return err
		}
		err = b.db.UpdateBidPaymentUrl(b.bid.Auction, b.bid.ExtraPrice, order.TrusteeShortlink.ShortURL)
		if err != nil {
			return err
		}
	} else {
		log.Warningf("[bid] Auction %v, extra_price %v | waiting for order %v to be processed", b.bid.Auction, b.bid.ExtraPrice, b.bid.PaymentID)
		market, err := marketpay.New()
		if err != nil {
			return err
		}
		paid, err := market.IsPaid(b.bid.PaymentID)
		if err != nil {
			return err
		}
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
				err := b.db.UpdateBidState(b.bid.Auction, b.bid.ExtraPrice, storage.BID_FAILED, err.Error())
				if err != nil {
					return err
				}
				b.bid.State = storage.BID_FAILED
				return err
			}
			receipt, err := helper.WaitReceipt(b.client, tx, 60)
			if err != nil {
				err := b.db.UpdateBidState(b.bid.Auction, b.bid.ExtraPrice, storage.BID_FAILED, err.Error())
				if err != nil {
					return err
				}
				b.bid.State = storage.BID_FAILED
				return err
			}
			if receipt.Status == 0 {
				err := b.db.UpdateBidState(b.bid.Auction, b.bid.ExtraPrice, storage.BID_FAILED, "receipt.Status == 0")
				if err != nil {
					return err
				}
				b.bid.State = storage.BID_FAILED
				return err
			}
			err = b.db.UpdateBidState(b.bid.Auction, b.bid.ExtraPrice, storage.BID_PAID, "")
			if err != nil {
				return err
			}
			b.bid.State = storage.BID_PAID
		}
	}
	return nil
}

func (b *BidMachine) processAccepted() error {
	err := b.db.UpdateBidState(b.bid.Auction, b.bid.ExtraPrice, storage.BID_PAYING, "")
	if err != nil {
		return err
	}
	b.bid.State = storage.BID_PAYING
	return nil
}
