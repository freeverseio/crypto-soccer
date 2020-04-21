package auctionmachine

import (
	"errors"

	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (m *AuctionMachine) ProcessPaying(market marketpay.IMarketPay) error {
	if m.auction.State != storage.AuctionPaying {
		return errors.New("Paying: wrong state")
	}

	// bid := bidmachine.FirstAlive(m.Bids)
	// if bid == nil {
	// 	m.auction.State = storage.AuctionFailed
	// 	m.auction.StateExtra = "Failed to pay"
	// 	return nil
	// }

	// bidMachine, err := bidmachine.New(
	// 	market,
	// 	&m.Auction,
	// 	bid,
	// 	m.contracts,
	// 	m.freeverse,
	// )
	// if err != nil {
	// 	return err
	// }

	// err = bidMachine.Process()
	// if err != nil {
	// 	return err
	// }
	// if bid.State == storage.BIDPAID {
	// 	log.Infof("[auction] %v PAYING -> PAID", m.Auction.UUID)
	// 	m.Auction.State = storage.AUCTION_PAID
	// }

	return nil
}
