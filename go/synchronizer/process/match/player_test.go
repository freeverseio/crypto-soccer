package match_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/process/match"
	"gotest.tools/assert"
)

func TestDefenceOfNullPlayer(t *testing.T) {
	t.Parallel()
	p := match.NewNullPlayer()
	defence, err := p.Defence(bc.Contracts.Assets)
	assert.NilError(t, err)
	assert.Equal(t, defence, uint16(0))
}

func TestDefenceOfPlayer(t *testing.T) {
	t.Parallel()
	p := match.NewPlayerFromSkills("14606253788909032162646379450304996475079674564248175")
	defence, err := p.Defence(bc.Contracts.Assets)
	assert.NilError(t, err)
	assert.Equal(t, defence, uint16(955))
}

func TestNewPlayer(t *testing.T) {
	t.Parallel()
	defence := uint16(50)
	speed := uint16(50)
	endurance := uint16(50)
	pass := uint16(50)
	shoot := uint16(50)
	dayOfBirthUnix := uint16(13344)
	generation := uint8(0)
	playerID := big.NewInt(2132321)
	potential := uint8(3)
	forwardness := uint8(0)
	leftishness := uint8(0)
	aggressiveness := uint8(0)
	alignedEndOfLastHalf := true
	redCardLastGame := false
	gamesNonStopping := uint8(0)
	injuryWeeksLeft := uint8(0)
	substitutedLastHalf := false
	p, err := match.NewPlayer(
		bc.Contracts,
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
	assert.Equal(t, p.Skills().String(), "730756529746917314243503421506698786561881762037810")
	value, err := p.Defence(bc.Contracts.Assets)
	assert.NilError(t, err)
	assert.Equal(t, value, defence)
	value, err = p.Speed(bc.Contracts.Assets)
	assert.NilError(t, err)
	assert.Equal(t, value, speed)
	birth, err := p.Birth(bc.Contracts.Assets)
	assert.NilError(t, err)
	assert.Equal(t, birth.Unix(), int64(dayOfBirthUnix)*3600*24)
}
