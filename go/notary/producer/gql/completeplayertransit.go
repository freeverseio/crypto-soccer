package gql

import (
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CompletePlayerTransit(args struct {
	Input input.CompletePlayerTransitInput
}) (graphql.ID, error) {
	log.Debugf("CompletePlayerTransit %v", args)

	id := graphql.ID(args.Input.PlayerId)

	return id, b.push(args.Input)
}
