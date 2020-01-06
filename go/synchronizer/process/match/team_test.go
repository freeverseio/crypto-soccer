package match_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/process/match"
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
	if err != nil {
		t.Fatal(err)
	}
	team.Players[0] = match.NewPlayer("4544")
	skills := team.Skills()
	if len(skills) != 25 {
		t.Fatalf("Wrong team skills size %v", len(skills))
	}
	if skills[0].Cmp(big.NewInt(4544)) != 0 {
		t.Fatalf("Wrong state %v", skills[0])
	}
}
