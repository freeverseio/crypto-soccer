package gql

import (
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CancelAllOffersBySeller(args struct {
	Input input.CancelAllOffersBySellerInput
}) (graphql.ID, error) {
	log.Debugf("CancelAllOffersBySeller %v", args)

	id := args.Input.PlayerId

	isOwner, err := args.Input.IsSignerOwnerOfPlayer(b.contracts)
	if err != nil {
		return id, err
	}
	if !isOwner {
		return id, fmt.Errorf("signer is not the owner of playerId %v", args.Input.PlayerId)
	}

	tx, err := b.service.Begin()
	if err != nil {
		return id, err
	}

	if err := tx.CancelAllOffersByPlayerId(string(id)); err != nil {
		tx.Rollback()
		return id, err
	}

	return id, tx.Commit()
}
