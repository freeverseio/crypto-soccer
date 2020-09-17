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
		return id, errors.New("empty OfferId when trying to cancel an offer")
	}

	offer, err := tx.Offer(string(args.Input.OfferId))
	if err != nil {
		return id, err
	}
	if offer == nil {
		return id, errors.New("could not find an offer to cancel")
	}

	if offer.State != storage.OfferStarted {
		return id, errors.New("cannot cancel an offer unless it is in Started state")
	}

	signer, err := args.Input.SignerAddress()
	if err != nil {
		return id, err
	}
	log.Warning("aaa")

	log.Warning(signer.Hex())
	log.Warning(offer.Seller)

	if (signer.Hex() != offer.Seller) && (signer.Hex() != offer.Buyer) {
		return id, errors.New("Signer of Canceloffer is neither the Seller nor the Buyer")
	}

	if err := tx.OfferCancel(string(args.Input.OfferId)); err != nil {
		tx.Rollback()
		return id, err
	}
	return id, tx.Commit()
}
