package match

import "math/big"

type Player struct {
	state *big.Int
}

func NewPlayer() *Player {
	var player Player
	player.state = big.NewInt(0)
	return &player
}
