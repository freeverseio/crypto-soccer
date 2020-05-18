package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) SetTeamName(args struct{ Input input.SetTeamNameInput }) (graphql.ID, error) {
	result := graphql.ID("")

	if b.ch == nil {
		return result, errors.New("internal error: no channel")
	}

	select {
	case b.ch <- args.Input:
	default:
		log.Warning("channel is full")
		return result, errors.New("channel is full")
	}

	return result, nil
}
