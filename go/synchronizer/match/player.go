package match

import "math/big"

type Player struct {
	state *big.Int
}

func NewPlayer(state string) *Player {
	var player Player
	player.state, _ = new(big.Int).SetString(state, 10)
	return &player
}

func NewPlayerDummy() *Player {
	var player Player
	player.state, _ = new(big.Int).SetString("3618502788679706519584493278137328010759678544985289844045583163109752700928", 10)
	return &player
}
