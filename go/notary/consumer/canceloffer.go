package consumer

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func CancelOffer(service storage.StorageService, tx *sql.Tx, in input.CancelOfferInput) error {
	offer, err := service.Offer(tx, string(in.OfferId))
	if err != nil {
		return err
	}
	if offer == nil {
		return errors.New("trying to cancel a nil offer")
	}
	if offer.State != storage.OfferStarted {
		return fmt.Errorf("not possible to cancel an offer in state %v", offer.State)
	}

	offer.State = storage.OfferCancelled
	return service.OfferUpdate(tx, offer)
}
