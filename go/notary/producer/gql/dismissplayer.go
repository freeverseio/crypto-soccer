package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) DismissPlayer(args struct{ Input input.DismissPlayerInput }) (graphql.ID, error) {
	log.Debugf("DismissPlayer %v", args)

	id := graphql.ID(args.Input.PlayerId)

	if b.ch == nil {
		return id, errors.New("internal error: no channel")
	}

	isValid, err := args.Input.VerifySignature()
	if err != nil {
		return graphql.ID(id), err
	}
	if !isValid {
		return graphql.ID(id), errors.New("Invalid signature")
	}

	isOwner, err := args.Input.IsSignerOwner(b.contracts)
	if err != nil {
		return id, err
	}
	if !isOwner {
		return id, errors.New("Not player owner")
	}

	select {
	case b.ch <- args.Input:
	default:
		log.Warning("channel is full")
		return id, errors.New("channel is full")
	}

	return id, nil
}
