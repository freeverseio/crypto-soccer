package consumer

import (
	"crypto/ecdsa"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/marketpay"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func ProcessAuctions(
	service storage.Tx,
	market marketpay.MarketPayService,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
) error {
	auctions, err := service.AuctionPendingAuctions()
	if err != nil {
		return err
	}

	for _, auction := range auctions {
		if err := processAuction(
			service,
			market,
			auction,
			pvc,
			contracts,
		); err != nil {
			log.Error(err)
		}
	}
	return nil
}

func processAuction(
	service storage.Tx,
	market marketpay.MarketPayService,
	auction storage.Auction,
	pvc *ecdsa.PrivateKey,
	contracts contracts.Contracts,
) error {
	bids, err := service.Bids(auction.ID)
	if err != nil {
		return err
	}

	am, err := auctionmachine.New(auction, bids, contracts, pvc)
	if err != nil {
		return err
	}
	if err := am.Process(market); err != nil {
		return err
	}
	if err := service.AuctionUpdate(am.Auction()); err != nil {
		return err
	}
	for _, bid := range am.Bids() {
		if err := service.BidUpdate(bid); err != nil {
			return err
		}
	}
	return nil
}
