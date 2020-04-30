package gql

import "github.com/graph-gophers/graphql-go"

type WorldPlayer struct {
	playerId          graphql.ID
	name              string
	dayOfBirth        int32
	preferredPosition string
	defence           int32
	speed             int32
	pass              int32
	shoot             int32
	endurance         int32
	potential         int32
	validUntil        string
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

func (b WorldPlayer) PlayerId() graphql.ID {
	return b.playerId
}

func (b WorldPlayer) Name() string {
	return b.name
}

func (b WorldPlayer) ValidUntil() string {
	return b.validUntil
}

func (b WorldPlayer) DayOfBirth() int32 {
	return b.dayOfBirth
}

func (b WorldPlayer) PreferredPosition() string {
	return b.preferredPosition
}

func (b WorldPlayer) Defence() int32 {
	return b.defence
}

func (b WorldPlayer) Speed() int32 {
	return b.speed
}

func (b WorldPlayer) Pass() int32 {
	return b.pass
}

func (b WorldPlayer) Shoot() int32 {
	return b.shoot
}

func (b WorldPlayer) Endurance() int32 {
	return b.endurance
}

func (b WorldPlayer) Potential() int32 {
	return b.potential
}
