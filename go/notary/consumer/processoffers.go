package consumer

import (
	"crypto/ecdsa"
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	log "github.com/sirupsen/logrus"
)

func ProcessOffers(
	market marketpay.IMarketPay,
	tx *sql.Tx,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
) error {
	service := postgres.NewOfferHistoryService(tx)
	offers, err := service.PendingOffers()
	if err != nil {
		return err
	}

	for _, offer := range offers {
		if err := processOffer(
			market,
			tx,
			offer,
			pvc,
			contracts,
		); err != nil {
			log.Error(err)
		}
	}
	return nil
}

func processOffer(
	market marketpay.IMarketPay,
	tx *sql.Tx,
	offer storage.Offer,
	pvc *ecdsa.PrivateKey,
	contracts contracts.Contracts,
) error {
	service := postgres.NewOfferHistoryService(tx)
	bid := CreateBidInput{
		Signature: offer.Signature,
		AuctionId: offer.AuctionID,
		ExtraPrice: 0,
		Rnd: offer.Rnd,
		TeamId: offer.TeamID,
		IsOffer: true
	}
	
	err := CreateBid(tx, bid)

	if err != nil {
		log.Error(err)
		return err
	}
	offer.state = storage.OfferAccepted 
	service.Update(offer)

	return nil
}
