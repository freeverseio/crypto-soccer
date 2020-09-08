package auctionmachine

import (
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/marketpay"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	log "github.com/sirupsen/logrus"
)

func (b *AuctionMachine) ProcessWithdrawableBySeller(market marketpay.MarketPayService) error {
	if err := b.checkState(storage.AuctionWithdrableBySeller); err != nil {
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
		b.SetState(storage.AuctionValidation, "")
	case "PUBLISHED":
		log.Debugf("Auction[%v|%v] order in state %v", b.auction.ID, b.auction.State, order.Status)
	default:
		log.Errorf("Auction[%v|%v] order in unknown state %v", b.auction.ID, b.auction.State, order.Status)
	}

	return nil
}
