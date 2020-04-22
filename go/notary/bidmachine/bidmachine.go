package bidmachine

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	log "github.com/sirupsen/logrus"
)

type BidMachine struct {
	market          marketpay.IMarketPay
	auction         storage.Auction
	bid             *storage.Bid
	contracts       contracts.Contracts
	freeverse       *ecdsa.PrivateKey
	PostAuctionTime int64
}

func New(
	market marketpay.IMarketPay,
	auction storage.Auction,
	bid *storage.Bid,
	contracts contracts.Contracts,
	freeverse *ecdsa.PrivateKey,
) (*BidMachine, error) {
	if market == nil {
		return nil, errors.New("No market instance given")
	}
	if auction.State != storage.AuctionPaying {
		return nil, errors.New("Auction is not in PAYING state")
	}
	if auction.ID != bid.AuctionID {
		return nil, errors.New("Bid of wrong auction")
	}
	PostAuctionTime, err := contracts.ConstantsGetters.GetPOSTAUCTIONTIME(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	return &BidMachine{
		market,
		auction,
		bid,
		contracts,
		freeverse,
		PostAuctionTime.Int64(),
	}, nil
}

func FirstAlive(bids []storage.Bid) *storage.Bid {
	// first searching for PAYING bid
	for i := range bids {
		if bids[i].State == storage.BidPaying {
			return &bids[i]
		}
	}
	// then search for the highest ACCEPTED bid
	idx := -1
	extraPrice := int64(-1)
	for i, bid := range bids {
		if bid.State == storage.BidAccepted {
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
	return &bids[idx]
}

func (b *BidMachine) Process() error {
	switch b.bid.State {
	case storage.BidPaying:
		return b.processPaying()
	case storage.BidAccepted:
		return b.processAccepted()
	case storage.BidFailed:
	default:
		return fmt.Errorf("Unknown bid state %v", b.bid.State)
	}
	return nil
}

func (b *BidMachine) processPaying() error {
	now := time.Now().Unix()
	if now > b.bid.PaymentDeadline {
		b.bid.State = storage.BidFailed
		b.bid.StateExtra = "Expired"
		return nil
	}

	log.Warningf("[bid] Auction %v, extra_price %v | waiting for order %v to be processed", b.bid.AuctionID, b.bid.ExtraPrice, b.bid.PaymentID)
	order, err := b.market.GetOrder(b.bid.PaymentID)
	if err != nil {
		return err
	}
	isPaid := b.market.IsPaid(*order)
	if isPaid {
		b.bid.State = storage.BidPaid
		b.bid.StateExtra = ""
	}

	return nil
}

func (b *BidMachine) processAccepted() error {
	log.Infof("[bid] Auction %v, extra_price %v | create MarketPay order", b.bid.AuctionID, b.bid.ExtraPrice)
	price := fmt.Sprintf("%.2f", float64(b.auction.Price+b.bid.ExtraPrice)/100.0)
	name := "Freeverse Player transaction"
	order, err := b.market.CreateOrder(name, price)
	if err != nil {
		b.bid.State = storage.BidFailed
		b.bid.StateExtra = ""
		return nil
	}
	b.bid.PaymentID = order.TrusteeShortlink.Hash
	b.bid.PaymentURL = order.TrusteeShortlink.ShortURL
	b.bid.PaymentDeadline = b.auction.ValidUntil + b.PostAuctionTime

	b.bid.State = storage.BidPaying
	return nil
}
