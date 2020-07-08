package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CompletePlayerTransit(args struct {
	Input input.CompletePlayerTransitInput
}) (graphql.ID, error) {
	log.Debugf("DismissPlayer %v", args)

	id := graphql.ID(args.Input.PlayerId)

	select {
	case b.ch <- args.Input:
	default:
		log.Warning("channel is full")
		return id, errors.New("channel is full")
	}

	return id, nil
}
