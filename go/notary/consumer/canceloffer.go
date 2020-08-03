package consumer

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func CancelOffer(service storage.StorageService, tx *sql.Tx, in input.CancelOfferInput) error {
	offerId, err := strconv.ParseInt(string(in.OfferId), 10, 64)
	if err != nil {
		return err
	}
	offer, err := service.Offer(tx, offerId)
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
	return service.OfferUpdate(tx, *offer)
}
