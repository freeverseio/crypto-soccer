package consumer

import (
	"crypto/ecdsa"
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
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
		bids := []storage.Bid{}
		am, err := auctionmachine.New(auction, bids, contracts, pvc)
		if err = am.Process(marketpay.New()); err != nil {

		}
	}
	return nil
}
