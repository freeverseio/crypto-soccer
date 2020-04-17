package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CreateBid(args struct{ Input input.CreateBidInput }) (graphql.ID, error) {
	log.Debugf("CreateBid %v", args)

	id := graphql.ID(args.Input.ID())

	if b.ch == nil {
		return id, errors.New("internal error: no channel")
	}

	isValid, err := args.Input.VerifySignature(b.contracts, *auction)
	if err != nil {
		return graphql.ID(id), err
	}
	if !isValid {
		return graphql.ID(id), errors.New("Invalid signature")
	}

	select {
	case b.ch <- args.Input:
	default:
		log.Warning("channel is full")
		return id, errors.New("channel is full")
	}
	return id, nil
}
