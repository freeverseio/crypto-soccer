package names

import (
	"math/big"

	"github.com/Pallinder/sillyname-go"
)

func GeneratePlayerName(playerId *big.Int) string {
	_ = playerId
	return sillyname.GenerateStupidName()
}

func GenerateTeamName(teamId *big.Int) string {
	_ = teamId
	return sillyname.GenerateStupidName()
}
