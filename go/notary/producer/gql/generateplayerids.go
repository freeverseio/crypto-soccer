package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) GeneratePlayerIDs(args struct{ Input input.GeneratePlayerIDsInput }) ([]graphql.ID, error) {
	log.Debugf("GeneratePlayerIDs %v", args)

	if b.ch == nil {
		return []graphql.ID{}, errors.New("internal error: no channel")
	}

	return []graphql.ID{}, errors.New("not implemented")
}
