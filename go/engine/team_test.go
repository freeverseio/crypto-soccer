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
	team, err := engine.NewTeam(bc.Contracts)
	if err != nil {
		t.Fatal(err)
	}
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
	team, err := engine.NewTeam(bc.Contracts)
	assert.NilError(t, err)
	skills := team.Skills()
	for _, skill := range skills {
		assert.Equal(t, skill.String(), "0")
	}
	team.Players[2] = engine.NewPlayerFromSkills("4544")
	skills = team.Skills()
	assert.Equal(t, skills[2].String(), "4544")
}
