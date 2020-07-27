package consumer

import (
	"crypto/ecdsa"
	"database/sql"
	"time"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	"github.com/graph-gophers/graphql-go"
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

	if offer.ValidUntil < time.Now().Unix() {
		offer.State = storage.OfferEnded
		offer.StateExtra = "Offer expired when processing"
		service.Update(offer)

		return nil
	}

	bid := input.CreateBidInput{
		Signature:  offer.Signature,
		AuctionId:  graphql.ID(offer.AuctionID),
		ExtraPrice: 0,
		Rnd:        int32(offer.Rnd),
		TeamId:     offer.TeamID,
	}

	err := CreateBid(tx, bid)

	if err != nil {
		log.Error(err)
		offer.State = storage.OfferFailed
		offer.StateExtra = "Could not create bid"
		service.Update(offer)
		return err
	}

	return nil
}
