package gql

import "github.com/graph-gophers/graphql-go"

type WorldPlayer struct {
	playerId graphql.ID
	name     string
}

func (b *WorldPlayer) PlayerId() graphql.ID {
	return b.playerId
}

func (b *WorldPlayer) Name() string {
	return b.name
}
