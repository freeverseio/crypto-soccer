package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CreateBid(args struct{ Input input.CreateBidInput }) (graphql.ID, error) {
	log.Infof("[notary|producer|gql] create bid %+v", args.Input)

	id, err := args.Input.ID(b.contracts)
	if err != nil {
		return graphql.ID(""), err
	}

	if b.ch == nil {
		return id, errors.New("internal error: no channel")
	}

	isValid, err := args.Input.VerifySignature(b.contracts)
	if err != nil {
		return graphql.ID(id), err
	}
	if !isValid {
		return graphql.ID(id), errors.New("Invalid signature")
	}

	return id, b.push(args.Input)
}
