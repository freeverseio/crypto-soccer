package gql

import "github.com/graph-gophers/graphql-go"

type WorldPlayer struct {
	playerId graphql.ID
	name     string
}

func NewWorldPlayer(
	playerId graphql.ID,
	name string,
) *WorldPlayer {
	player := WorldPlayer{}
	player.playerId = playerId
	player.name = name
	return &player
}

func (b *WorldPlayer) PlayerId() graphql.ID {
	return b.playerId
}

func (b *WorldPlayer) Name() string {
	return b.name
}
