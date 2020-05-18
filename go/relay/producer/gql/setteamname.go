package gql

import (
	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
)

func (b *Resolver) SetTeamName(args struct{ Input input.SetTeamNameInput }) (graphql.ID, error) {
	result := graphql.ID("")

	return result, nil
}
