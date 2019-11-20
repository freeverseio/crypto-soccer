package processor

import (
	"crypto/ecdsa"

	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"

	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"

	"github.com/freeverseio/crypto-soccer/go/contracts/market"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

type Processor struct {
	db        *storage.Storage
	client    *ethclient.Client
	assets    *market.Market
	freeverse *ecdsa.PrivateKey
}

func NewProcessor(db *storage.Storage, ethereumClient *ethclient.Client, assetsContract *market.Market, freeverse *ecdsa.PrivateKey) (*Processor, error) {
	return &Processor{db, ethereumClient, assetsContract, freeverse}, nil
}

func (b *Processor) Process() error {
	log.Info("Processing")

	openedAuctions, err := b.db.GetOpenAuctions()
	if err != nil {
		return err
	}

	for _, auction := range openedAuctions {
		bids, err := b.db.GetBidsOfAuction(auction.UUID)
		if err != nil {
			return err
		}

		machine, err := auctionmachine.New(
			auction,
			bids,
			b.assets,
			b.freeverse,
			b.client,
		)
		if err != nil {
			return err
		}
		err = machine.Process()
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
	return b.db.UpdateAuctionState(auction.UUID, auction.State)
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
	}
	return nil
}
