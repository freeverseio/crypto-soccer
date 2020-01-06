package match_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/process/match"
	"gotest.tools/assert"
)

func TestTeamStateDefault(t *testing.T) {
	team, err := match.NewTeam(bc.Contracts)
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
}

func TestTeamSkills(t *testing.T) {
	team, err := match.NewTeam(bc.Contracts)
	assert.NilError(t, err)
	skills := team.Skills()
	for _, skill := range skills {
		assert.Equal(t, skill.String(), "0")
	}
	team.Players[2] = match.NewPlayer("4544")
	skills = team.Skills()
	assert.Equal(t, skills[2].String(), "4544")
}
