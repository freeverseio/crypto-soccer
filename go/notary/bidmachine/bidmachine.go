package bidmachine

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/helper"
	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	log "github.com/sirupsen/logrus"
)

type BidMachine struct {
	auction         *storage.Auction
	bid             *storage.Bid
	contracts       *contracts.Contracts
	freeverse       *ecdsa.PrivateKey
	signer          *signer.Signer
	postAuctionTime *big.Int
}

func New(
	auction *storage.Auction,
	bid *storage.Bid,
	contracts *contracts.Contracts,
	freeverse *ecdsa.PrivateKey,
) (*BidMachine, error) {
	if auction.State != storage.AUCTION_PAYING {
		return nil, errors.New("Auction is not in PAYING state")
	}
	if auction.UUID != bid.Auction {
		return nil, errors.New("Bid of wrong auction")
	}
	postAuctionTime, err := contracts.Market.POSTAUCTIONTIME(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	return &BidMachine{
		auction,
		bid,
		contracts,
		freeverse,
		signer.NewSigner(contracts, freeverse),
		postAuctionTime,
	}, nil
}

func FirstAlive(bids []*storage.Bid) *storage.Bid {
	// first searching for PAYING bid
	for i := range bids {
		if bids[i].State == storage.BIDPAYING {
			return bids[i]
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
	if idx == -1 {
		return nil
	}
	return bids[idx]
}

func (b *BidMachine) Process() error {
	switch b.bid.State {
	case storage.BIDPAYING:
		return b.processPaying()
	case storage.BIDACCEPTED:
		return b.processAccepted()
	case storage.BIDFAILED:
		return nil
	default:
		return errors.New("Unknown bid state")
	}
}

func (b *BidMachine) processPaying() error {
	now := time.Now().Unix()
	if now > b.bid.PaymentDeadline.Int64() {
		b.bid.State = storage.BIDFAILED
		b.bid.StateExtra = "Expired"
		return nil
	}
	if b.bid.PaymentID == "" { // create order
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
	} else { // check if order is paid
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
			tx, err := b.contracts.Market.CompletePlayerAuction(
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
			receipt, err := helper.WaitReceipt(b.contracts.Client, tx, 60)
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
	if b.auction.ValidUntil == nil {
		return errors.New("nil valid until")
	}
	b.bid.State = storage.BIDPAYING
	b.bid.PaymentDeadline = new(big.Int).Add(b.auction.ValidUntil, b.postAuctionTime)
	return nil
}
