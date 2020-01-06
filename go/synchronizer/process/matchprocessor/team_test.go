package matchprocessor_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/process/matchprocessor"
)

func TestTeamStateDefault(t *testing.T) {
	team, err := matchprocessor.NewTeam(bc.Contracts)
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

func TestTeamState(t *testing.T) {
	team, err := matchprocessor.NewTeam(bc.Contracts)
	if err != nil {
		t.Fatal(err)
	}
	team.Players[0] = matchprocessor.NewPlayer("4544")
	states := team.Skills()
	if len(states) != 25 {
		t.Fatalf("Wrong states size %v", len(states))
	}
	if states[0].Cmp(big.NewInt(4544)) != 0 {
		t.Fatalf("Wrong state %v", states[0])
	}
}
