package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CreateAuction(args struct{ Input input.CreateAuctionInput }) (graphql.ID, error) {
	if b.c == nil {
		return graphql.ID(""), errors.New("internal error: no channel")
	}

	isValid, err := args.Input.VerifySignature()
	if err != nil {
		return graphql.ID(""), err
	}
	if !isValid {
		return graphql.ID(""), errors.New("Invalid signature")
	}

	select {
	case b.c <- args.Input:
	default:
		log.Warning("channel is full")
		return graphql.ID("ciao"), errors.New("channel is full")
	}

	return graphql.ID(args.Input.Signature), nil
}
