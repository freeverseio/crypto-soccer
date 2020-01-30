package engine_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/engine"
	"gotest.tools/assert"
)

func TestNullPlayer(t *testing.T) {
	t.Parallel()
	p := engine.NewNullPlayer()
	assert.Equal(t, p.DumpState(), "skills: 0")
}

func TestDefenceOfNullPlayer(t *testing.T) {
	t.Parallel()
	p := engine.NewNullPlayer()
	assert.Equal(t, p.Defence(), uint16(0))
}

func TestPlayerNewPlayerFromSkills(t *testing.T) {
	t.Parallel()
	p, err := engine.NewPlayerFromSkills(*bc.Contracts, "14606253788909032162646379450304996475079674564248175")
	assert.NilError(t, err)
	assert.Equal(t, p.Defence(), uint16(955))
	assert.Equal(t, p.Speed(), uint16(955))
	assert.Equal(t, p.Pass(), uint16(955))
	assert.Equal(t, p.Endurance(), uint16(955))
	assert.Equal(t, p.Potential(), uint16(955))
	assert.Equal(t, p.Shoot(), uint16(955))
	assert.Equal(t, p.DayOfBirth(), uint16(955))
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
	p, err := engine.NewPlayer(
		*bc.Contracts,
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
	assert.Equal(t, p.Defence(), defence)
	value, err := p.Speed(bc.Contracts.Assets)
	assert.NilError(t, err)
	assert.Equal(t, value, speed)
	value, err = p.BirthDayUnix(bc.Contracts.Assets)
	assert.NilError(t, err)
	assert.Equal(t, value, dayOfBirthUnix)
}
