package names

import (
	"github.com/Pallinder/sillyname-go"
	"math/big"
)

func GeneratePlayerName(playerId *big.Int) string {
	_ = playerId
	return sillyname.GenerateStupidName()
}
