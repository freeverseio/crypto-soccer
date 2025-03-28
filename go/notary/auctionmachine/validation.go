package auctionmachine

import (
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/marketpay"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func (b *AuctionMachine) ProcessValidation(market marketpay.MarketPayService) error {
	if err := b.checkState(storage.AuctionValidation); err != nil {
		return err
	}

	paidBids := storage.FindBids(b.bids, storage.BidPaid)
	if len(paidBids) != 1 {
		return fmt.Errorf("Paid bids should be 1 but found %v", len(paidBids))
	}
	paidBid := paidBids[0]

	order, err := market.GetOrder(paidBid.PaymentID)
	if err != nil {
		return err
	}

	switch order.Status {
	case "PENDING_VALIDATE":
		result, err := market.ValidateOrder(paidBid.PaymentID)
		if err != nil {
			return err
		}
		log.Infof("auction[%v|%v] validation result %v", b.auction.ID, b.auction.State, result)
		b.SetState(storage.AuctionValidation, result)
	case "PENDING_RELEASE":
		log.Infof("auction[%v|%v] pending release", b.auction.ID, b.auction.State)
	case "RELEASED":
		b.SetState(storage.AuctionEnded, "")
	default:
		log.Errorf("auction[%v|%v] order in unknown state %v", b.auction.ID, b.auction.State, order.Status)
	}

	return nil
}
