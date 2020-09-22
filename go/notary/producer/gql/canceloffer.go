package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CancelOffer(args struct{ Input input.CancelOfferInput }) (graphql.ID, error) {
	log.Debugf("CancelOffer %v", args)

	id := args.Input.OfferId

	tx, err := b.service.Begin()
	if err != nil {
		return id, err
	}

	if string(args.Input.OfferId) == "" {
		tx.Rollback()
		return id, errors.New("empty OfferId when trying to cancel an offer")
	}

	offer, err := tx.Offer(string(args.Input.OfferId))
	log.Warningf("Trying to look for offerid: %v\n", args.Input.OfferId)
	log.Warningf("Offer: %v\n", offer)
	if err != nil {
		tx.Rollback()
		return id, err
	}
	if offer == nil {
		tx.Rollback()
		return id, errors.New("could not find an offer to cancel")
	}

	log.Warningf("storage.OfferStarted: %v\n", storage.OfferStarted)
	log.Warningf("offer.state: %v\n", offer.State)

	if offer.State != storage.OfferStarted {
		tx.Rollback()
		return id, errors.New("cannot cancel an offer unless it is in Started state")
	}

	signer, err := args.Input.SignerAddress()
	if err != nil {
		tx.Rollback()
		return id, err
	}

	if (signer.Hex() != offer.Seller) && (signer.Hex() != offer.Buyer) {
		tx.Rollback()
		return id, errors.New("Signer of Canceloffer is neither the Seller nor the Buyer")
	}

	if err := tx.OfferCancel(string(args.Input.OfferId)); err != nil {
		tx.Rollback()
		return id, err
	}
	return id, tx.Commit()
}
