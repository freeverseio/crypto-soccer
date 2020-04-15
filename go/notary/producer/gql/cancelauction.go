package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CancelAuction(args struct{ Input input.CancelAuctionInput }) (graphql.ID, error) {
	log.Debugf("CancelAuction %v", args)

	id := args.Input.ID

	if b.ch == nil {
		return graphql.ID(id), errors.New("internal error: no channel")
	}

	isValid, err := args.Input.VerifySignature()
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
		return graphql.ID(id), errors.New("channel is full")
	}
	return graphql.ID(id), nil
}
