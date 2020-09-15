package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CancelOffer(args struct{ Input input.CancelOfferInput }) (graphql.ID, error) {
	log.Debugf("CancelOffer %v", args)

	id := args.Input.OfferId

	if b.ch == nil {
		return id, errors.New("internal error: no channel")
	}

	// isValid, err := args.Input.VerifySignature()
	// if err != nil {
	// 	return id, err
	// }
	// if !isValid {
	// 	return id, errors.New("Invalid signature")
	// }

	tx, err := b.service.Begin()
	if err != nil {
		return id, err
	}
	if err := tx.OfferCancel(string(args.Input.OfferId)); err != nil {
		tx.Rollback()
		return id, err
	}

	return id, tx.Commit()
}
