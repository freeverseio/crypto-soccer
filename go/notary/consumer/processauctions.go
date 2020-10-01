package consumer

import (
	"crypto/ecdsa"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/marketpay"
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
