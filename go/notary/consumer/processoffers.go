package consumer

import (
	"database/sql"
	"time"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func ProcessOffers(
	service storage.StorageService,
	tx *sql.Tx,
) error {
	offers, err := service.OfferPendingOffers(tx)
	if err != nil {
		return err
	}

	for _, offer := range offers {
		if err := processOffer(
			service,
			tx,
			offer,
		); err != nil {
			log.Error(err)
		}
	}
	return nil
}

func processOffer(
	service storage.StorageService,
	tx *sql.Tx,
	offer storage.Offer,
) error {
	if offer.ValidUntil < time.Now().Unix() {
		offer.State = storage.OfferEnded
		offer.StateExtra = "Offer expired"
		if err := service.OfferUpdate(tx, offer); err != nil {
			return err
		}
	}

	return nil
}
