package consumer

import (
	"crypto/ecdsa"
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func ProcessAuctions(
	tx *sql.Tx,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
) error {
	auctions, err := storage.PendingAuctions(tx)
	if err != nil {
		return err
	}

	for _, auction := range auctions {
		if err := processAuction(tx, auction, pvc, contracts); err != nil {
			log.Error(err)
		}
	}
	return nil
}

func processAuction(
	tx *sql.Tx,
	auction storage.Auction,
	pvc *ecdsa.PrivateKey,
	contracts contracts.Contracts,
) error {
	bids, err := storage.BidsByAuctionID(tx, auction.ID)
	if err != nil {
		return err
	}
	am, err := auctionmachine.New(auction, bids, contracts, pvc)
	if err != nil {
		return err
	}
	if err := am.Process(marketpay.New()); err != nil {
		return err
	}
	if err := am.Auction().Update(tx); err != nil {
		return err
	}
	return nil
}
