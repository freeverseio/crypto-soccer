package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) DismissPlayer(args struct{ Input input.DismissPlayerInput }) (graphql.ID, error) {
	log.Infof("[notary|producer|gql] dismiss player %+v", args.Input)

	id := graphql.ID(args.Input.PlayerId)

	if b.ch == nil {
		return id, errors.New("internal error: no channel")
	}

	// isValid, err := args.Input.VerifySignature()
	// if err != nil {
	// 	return graphql.ID(id), err
	// }
	// if !isValid {
	// 	return graphql.ID(id), errors.New("invalid signature")
	// }

	isOwner, err := args.Input.IsSignerOwner(b.contracts)
	if err != nil {
		return id, err
	}
	if !isOwner {
		return id, errors.New("not player owner")
	}

	return id, b.push(args.Input)
}
