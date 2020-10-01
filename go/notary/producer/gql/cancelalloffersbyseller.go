package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
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
	log.Debugf("Input signature: %v\n", args.Input.Signature)
	log.Debugf("Input playerId: %v\n", args.Input.PlayerId)
	existingOffers, err := tx.OffersByPlayerId(string(args.Input.PlayerId))
	if err != nil {
		tx.Rollback()
		return id, errors.New("could not find existing offers")
	}
	if existingOffers == nil {
		tx.Rollback()
		return id, errors.New("could not find an existing offers to cancel")
	}

	for _, offer := range existingOffers {
		log.Debugf("Trying to cancel offer: %v\n", offer)

		if offer.State != storage.OfferStarted {
			tx.Rollback()
			return id, errors.New("cannot cancel an offer unless it is in Started state")
		}

		signer, err := args.Input.SignerAddress()
		if err != nil {
			tx.Rollback()
			return id, err
		}
		log.Debugf("Signer from input: %v\n", signer.Hex())
		log.Debugf("Offer seller: %v\n", offer.Seller)
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
