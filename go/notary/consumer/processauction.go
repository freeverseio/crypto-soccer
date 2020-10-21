package consumer

import (
	"crypto/ecdsa"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/marketpay"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func ProcessAuction(
	service storage.Tx,
	market marketpay.MarketPayService,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	auctionID string,
	shouldQueryMarketPay bool,
) error {
	auction, err := service.Auction(auctionID)
	if err != nil {
		return err
	}
	if err := processAuction(
		service,
		market,
		*auction,
		pvc,
		contracts,
		shouldQueryMarketPay,
	); err != nil {
		log.Error(err)
	}
	return nil
}

func processAuction(
	service storage.Tx,
	market marketpay.MarketPayService,
	auction storage.Auction,
	pvc *ecdsa.PrivateKey,
	contracts contracts.Contracts,
	shouldQueryMarketPay bool,
) error {
	bids, err := service.Bids(auction.ID)
	if err != nil {
		return err
	}

	am, err := auctionmachine.New(auction, bids, contracts, pvc, shouldQueryMarketPay)
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
