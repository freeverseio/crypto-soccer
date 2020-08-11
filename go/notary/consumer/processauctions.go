package consumer

import (
	"crypto/ecdsa"
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/marketpay"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func ProcessAuctions(
	service storage.StorageService,
	market marketpay.MarketPayService,
	tx *sql.Tx,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
) error {
	auctions, err := service.AuctionPendingAuctions(tx)
	if err != nil {
		return err
	}

	for _, auction := range auctions {
		if err := processAuction(
			service,
			market,
			tx,
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
	service storage.StorageService,
	market marketpay.MarketPayService,
	tx *sql.Tx,
	auction storage.Auction,
	pvc *ecdsa.PrivateKey,
	contracts contracts.Contracts,
) error {
	bids, err := service.Bids(tx, auction.ID)
	if err != nil {
		return err
	}
	offer, err := service.OfferByAuctionId(tx, auction.ID)
	if err != nil {
		return err
	}
	am, err := auctionmachine.New(auction, bids, *offer, contracts, pvc)
	if err != nil {
		return err
	}
	if err := am.Process(market); err != nil {
		return err
	}
	if err := service.AuctionUpdate(tx, am.Auction()); err != nil {
		return err
	}
	for _, bid := range am.Bids() {
		if err := service.BidUpdate(tx, bid); err != nil {
			return err
		}
	}
	return nil
}
