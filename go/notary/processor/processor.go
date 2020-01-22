package processor

import (
	"crypto/ecdsa"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	log "github.com/sirupsen/logrus"
)

type Processor struct {
	db        *storage.Storage
	contracts *contracts.Contracts
	freeverse *ecdsa.PrivateKey
	market    *marketpay.MarketPay
}

func NewProcessor(db *storage.Storage, contracts *contracts.Contracts, freeverse *ecdsa.PrivateKey) (*Processor, error) {
	market, err := marketpay.New()
	if err != nil {
		return nil, err
	}
	return &Processor{db, contracts, freeverse, market}, nil
}

func (b *Processor) Process() error {
	log.Info("Processing")

	openedAuctions, err := b.db.GetOpenAuctions()
	if err != nil {
		return err
	}

	for _, auction := range openedAuctions {
		log.Infof("[processor] process auction %v, state: %v", auction.UUID, string(auction.State))
		bids, err := b.db.GetBidsOfAuction(auction.UUID)
		if err != nil {
			return err
		}

		machine, err := auctionmachine.New(
			auction,
			bids,
			b.contracts,
			b.freeverse,
		)
		if err != nil {
			return err
		}
		err = machine.Process(b.market)
		if err != nil {
			return err
		}

		err = b.updateAuction(auction)
		if err != nil {
			return err
		}
		err = b.updateBids(bids)
		if err != nil {
			return err
		}

	}

	return nil
}

func (b *Processor) updateAuction(auction *storage.Auction) error {
	err := b.db.UpdateAuctionState(auction.UUID, auction.State, auction.StateExtra)
	if err != nil {
		return err
	}
	return b.db.UpdateAuctionPaymentUrl(auction.UUID, auction.PaymentURL)
}

func (b *Processor) updateBids(bids []*storage.Bid) error {
	for _, bid := range bids {
		err := b.db.UpdateBidPaymentID(bid.Auction, bid.ExtraPrice, bid.PaymentID)
		if err != nil {
			return err
		}
		err = b.db.UpdateBidPaymentUrl(bid.Auction, bid.ExtraPrice, bid.PaymentURL)
		if err != nil {
			return err
		}
		err = b.db.UpdateBidState(bid.Auction, bid.ExtraPrice, bid.State, bid.StateExtra)
		if err != nil {
			return err
		}
		err = b.db.UpdateBidPaymentDeadline(bid.Auction, bid.ExtraPrice, bid.PaymentDeadline)
		if err != nil {
			return err
		}
	}
	return nil
}
