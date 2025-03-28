package bidmachine

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/marketpay"
	log "github.com/sirupsen/logrus"
)

type BidMachine struct {
	market               marketpay.MarketPayService
	auction              storage.Auction
	bid                  *storage.Bid
	contracts            contracts.Contracts
	freeverse            *ecdsa.PrivateKey
	PostAuctionTime      int64
	shouldQueryMarketPay bool
}

func New(
	market marketpay.MarketPayService,
	auction storage.Auction,
	bid *storage.Bid,
	contracts contracts.Contracts,
	freeverse *ecdsa.PrivateKey,
	shouldQueryMarketPay bool,
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
		shouldQueryMarketPay,
	}, nil
}

func FirstAlive(bids []storage.Bid) *storage.Bid {
	// first searching for PAYING bid
	payingBids := storage.FindBids(bids, storage.BidPaying)
	if len(payingBids) != 0 {
		return payingBids[0]
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
		return nil
	default:
		return fmt.Errorf("Unknown bid state %v", b.bid.State)
	}
	return nil
}

func (b *BidMachine) processPaying() error {
	now := time.Now().Unix()
	if now > b.bid.PaymentDeadline {
		b.setState(storage.BidFailed, "expired")
		return nil
	}
	log.Debugf("[bid] Auction %v, extra_price %v | waiting for order %v to be processed", b.bid.AuctionID, b.bid.ExtraPrice, b.bid.PaymentID)

	if b.shouldQueryMarketPay {
		order, err := b.market.GetOrder(b.bid.PaymentID)
		log.Debugf("Order received from marketpay: %v\n", order)
		log.Debugf("error received from marketpay.GetOrder: %v\n", err)
		if err != nil {
			return err
		}
		isPaid := b.market.IsPaid(*order)
		if isPaid {
			b.setState(storage.BidPaid, "")
		}
		log.Debugf("[bid] Auction %v, extra_price %v | Order: %v, state: %v, is paid: %v | Bid New State: %v", b.bid.AuctionID, b.bid.ExtraPrice, b.bid.PaymentID, order.Status, isPaid, b.bid.State)
	}

	return nil
}

func (b *BidMachine) processAccepted() error {
	log.Infof("[bid] Auction %v extra_price %v create MarketPay order", b.bid.AuctionID, b.bid.ExtraPrice)
	price := fmt.Sprintf("%.2f", float64(b.auction.Price+b.bid.ExtraPrice)/100.0)
	name := "Freeverse Player transaction ID: " + b.auction.ID + ", price: " + price
	order, err := b.market.CreateOrder(name, price)
	if err != nil {
		b.setState(storage.BidFailed, err.Error())
		return nil
	}
	b.bid.PaymentID = order.TrusteeShortlink.Hash
	b.bid.PaymentURL = order.TrusteeShortlink.ShortURL
	b.bid.PaymentDeadline = b.auction.ValidUntil + b.PostAuctionTime

	b.setState(storage.BidPaying, "")
	return nil
}

func (b *BidMachine) setState(state storage.BidState, extra string) {
	if state == storage.BidFailed {
		log.Warnf("[bid] auction %v extra price %v in state %v with %v", b.auction.ID, b.bid.ExtraPrice, state, extra)
	}
	b.bid.State = state
	b.bid.StateExtra = extra
}
