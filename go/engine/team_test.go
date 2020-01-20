package engine_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/engine"
	"gotest.tools/assert"
	"gotest.tools/golden"
)

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

func TestTrainingPoints(t *testing.T) {
	cases := []struct {
		MatchLog       string
		ExpectedPoints uint64
	}{
		{"0", 0},
		{"828212031530063714069492904776115492597195551273105492225696825706808722", 15},
		{"993853943853037244967927045470764103456022166605194769473036213412693666", 18},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("log:%v expectedPoints:%v", tc.MatchLog, tc.ExpectedPoints), func(t *testing.T) {
			team, err := engine.NewTeam(bc.Contracts)
			assert.NilError(t, err)
			engineLog, _ := new(big.Int).SetString(tc.MatchLog, 10)
			err = team.SetTrainingPoints(bc.Contracts, engineLog)
			assert.NilError(t, err)
			assert.Equal(t, team.TrainingPoints, tc.ExpectedPoints)
		})
	}
}
