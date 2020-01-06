package match

import "math/big"

type Player struct {
	skills *big.Int
}

func NewPlayer(skills string) *Player {
	var player Player
	player.skills, _ = new(big.Int).SetString(skills, 10)
	return &player
}
