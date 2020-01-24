package engine_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/engine"
	"gotest.tools/assert"
	"gotest.tools/golden"
)

func TestNoOutOfGamePlayer(t *testing.T) {
	t.Parallel()
	result, err := bc.Contracts.Engineprecomp.NOOUTOFGAMEPLAYER(&bind.CallOpts{})
	assert.NilError(t, err)
	assert.Equal(t, result, engine.NoOutOfGamePlayer)
}

func TestHardInjury(t *testing.T) {
	t.Parallel()
	result, err := bc.Contracts.Engineprecomp.HARDINJURY(&bind.CallOpts{})
	assert.NilError(t, err)
	assert.Equal(t, result, engine.HardInjury)
}

func TestRedCard(t *testing.T) {
	t.Parallel()
	result, err := bc.Contracts.Engineprecomp.REDCARD(&bind.CallOpts{})
	assert.NilError(t, err)
	assert.Equal(t, result, engine.RedCard)
}

func TestSoftInjury(t *testing.T) {
	t.Parallel()
	result, err := bc.Contracts.Engineprecomp.SOFTINJURY(&bind.CallOpts{})
	assert.NilError(t, err)
	assert.Equal(t, result, engine.SoftInjury)
}

func TestTeamStateDefault(t *testing.T) {
	t.Parallel()
	team := engine.NewTeam(bc.Contracts)
	states := team.Skills()
	if len(states) != 25 {
		t.Fatalf("Wrong states size %v", len(states))
	}
	for _, state := range states {
		if state.Cmp(big.NewInt(0)) != 0 {
			t.Fatalf("Wrong %v", state)
		}
	}
	golden.Assert(t, team.DumpState(), t.Name()+".golden")
}

func TestTeamSkills(t *testing.T) {
	t.Parallel()
	team := engine.NewTeam(bc.Contracts)
	skills := team.Skills()
	for _, skill := range skills {
		assert.Equal(t, skill.String(), "0")
	}
	team.Players[2] = engine.NewPlayerFromSkills("4544")
	skills = team.Skills()
	assert.Equal(t, skills[2].String(), "4544")
}

func TestDefaultTactics(t *testing.T) {
	t.Parallel()
	substitutions := [3]uint8{11, 11, 11}
	substitutionsMinute := [3]uint8{2, 3, 4}
	formation := [14]uint8{0, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 25, 25, 25}
	extraAttack := [10]bool{false, false, false, false, false, false, false, false, false, false}
	tacticID := uint8(1)
	tactic, err := bc.Contracts.Engine.EncodeTactics(
		&bind.CallOpts{},
		substitutions,
		substitutionsMinute,
		formation,
		extraAttack,
		tacticID,
	)
	assert.NilError(t, err)
	assert.Equal(t, tactic.String(), engine.DefaultTactic().String())
}
