package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CancelAllOffersBySeller(args struct {
	Input input.CancelAllOffersBySellerInput
}) (graphql.ID, error) {
	log.Debugf("CancelAllOffersBySeller %v", args)

	id := args.Input.PlayerId

	tx, err := b.service.Begin()
	if err != nil {
		return id, err
	}

	existingStartedOffers, err := tx.OffersStartedByPlayerId(string(args.Input.PlayerId))
	if err != nil {
		tx.Rollback()
		return id, errors.New("could not find existing offers")
	}
	if existingStartedOffers == nil {
		tx.Rollback()
		return id, errors.New("could not find an existing offers to cancel")
	}

	for _, offer := range existingStartedOffers {

		signer, err := args.Input.SignerAddress()
		if err != nil {
			tx.Rollback()
			return id, err
		}

		if signer.Hex() != offer.Seller {
			tx.Rollback()
			return id, errors.New("Signer of Cancelalloffersbyseller is not the Seller")
		}

		if err := tx.OfferCancel(offer.AuctionID); err != nil {
			tx.Rollback()
			return id, err
		}
	}

	return id, tx.Commit()
}
