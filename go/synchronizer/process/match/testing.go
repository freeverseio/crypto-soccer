package match

import (
	"math/big"
	"testing"
	"time"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"gotest.tools/assert"
)

func CreateDummyPlayer(
	t *testing.T,
	contracts *contracts.Contracts,
	age uint16,
	defence uint16,
	speed uint16,
	endurance uint16,
	pass uint16,
	shoot uint16,
) *Player {
	nowDays := time.Now().Unix() / 3600 / 24
	dayOfBirthUnix := uint16(nowDays-int64(age)*365/7) - 1
	generation := uint8(0)
	playerID := big.NewInt(2132321)
	potential := uint8(3)
	forwardness := uint8(3)
	leftishness := uint8(7)
	aggressiveness := uint8(0)
	alignedEndOfLastHalf := true
	redCardLastGame := false
	gamesNonStopping := uint8(0)
	injuryWeeksLeft := uint8(0)
	substitutedLastHalf := false
	p, err := NewPlayer(
		contracts,
		playerID,
		defence,
		speed,
		endurance,
		pass,
		shoot,
		dayOfBirthUnix,
		generation,
		potential,
		forwardness,
		leftishness,
		aggressiveness,
		alignedEndOfLastHalf,
		redCardLastGame,
		gamesNonStopping,
		injuryWeeksLeft,
		substitutedLastHalf,
	)
	assert.NilError(t, err)
	return p
}
