package match

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"gotest.tools/assert"
)

func CreateDummyPlayer(
	t *testing.T,
	contracts *contracts.Contracts,
	defence uint16,
	speed uint16,
	endurance uint16,
	pass uint16,
	shoot uint16,
) *Player {
	dayOfBirth := uint16(17078)
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
		dayOfBirth,
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
