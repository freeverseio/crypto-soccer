package auctionmachine

import (
	"fmt"

	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	log "github.com/sirupsen/logrus"
)

func (b *AuctionMachine) ProcessWithdrawableBySeller(market marketpay.IMarketPay) error {
	if b.State() != storage.AuctionWithdrableBySeller {
		return fmt.Errorf("Wrong state %v", b.State())
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

	log.Warningf("TODO auction %v in state %v order status is %v", b.auction.ID, b.auction.State, order.Status)

	return nil
}
