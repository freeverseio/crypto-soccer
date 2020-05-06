package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) SubmitPlayerPurchase(args struct {
	Input input.SubmitPlayerPurchaseInput
}) (graphql.ID, error) {
	log.Debugf("GeneratePlayerIDs %v", args)

	result := graphql.ID(args.Input.PlayerId)

	isValid, err := args.Input.IsValidSignature()
	if err != nil {
		return result, err
	}
	if !isValid {
		return result, errors.New("Invalid signature")
	}

	isOwner, err := args.Input.IsSignerOwner(b.contracts)
	if err != nil {
		return result, err
	}
	if !isOwner {
		return result, errors.New("Not team owner")
	}

	return result, errors.New("not implemented")
}
