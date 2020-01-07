package match_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/process/match"
	"gotest.tools/assert"
)

func TestTeamStateDefault(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	team, err := match.NewTeam(bc.Contracts)
	assert.NilError(t, err)
	skills := team.Skills()
	for _, skill := range skills {
		assert.Equal(t, skill.String(), "0")
	}
	team.Players[2] = match.NewPlayerFromSkills("4544")
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
			team, err := match.NewTeam(bc.Contracts)
			assert.NilError(t, err)
			matchLog, _ := new(big.Int).SetString(tc.MatchLog, 10)
			err = team.SetTrainingPoints(bc.Contracts, matchLog)
			assert.NilError(t, err)
			assert.Equal(t, team.TrainingPoints, tc.ExpectedPoints)
		})
	}
}
