package consumer

import (
	"crypto/ecdsa"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/marketpay"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func ProcessAuctionFromWebhook(
	service storage.Tx,
	market marketpay.MarketPayService,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	shouldQueryMarketPay bool,
	auctionID string,
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
