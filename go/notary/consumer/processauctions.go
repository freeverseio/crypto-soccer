package consumer

import (
	"crypto/ecdsa"
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func ProcessAuctions(
	market marketpay.IMarketPay,
	tx *sql.Tx,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
) error {
	service := postgres.NewAuctionHistoryService(tx)
	auctions, err := service.PendingAuctions()
	if err != nil {
		return err
	}

	for _, auction := range auctions {
		if err := processAuction(
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
	market marketpay.IMarketPay,
	tx *sql.Tx,
	auction storage.Auction,
	pvc *ecdsa.PrivateKey,
	contracts contracts.Contracts,
) error {
	service := postgres.NewAuctionHistoryService(tx)
	bids, err := service.Bid().Bids(auction.ID)
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
	if err := service.Update(am.Auction()); err != nil {
		return err
	}
	for _, bid := range am.Bids() {
		if err := service.Bid().Update(bid); err != nil {
			return err
		}
	}
	return nil
}
