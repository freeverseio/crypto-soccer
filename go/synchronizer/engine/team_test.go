package engine_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/engine"
	"gotest.tools/assert"
	"gotest.tools/golden"
)

func TestTeamStateDefault(t *testing.T) {
	t.Parallel()
	team := engine.NewTeam()
	states := team.Skills()
	if len(states) != 25 {
		t.Fatalf("Wrong states size %v", len(states))
	}
	for _, state := range states {
		if state.Cmp(big.NewInt(0)) != 0 {
			t.Fatalf("Wrong %v", state)
		}
	}
	golden.Assert(t, dump.Sdump(team), t.Name()+".golden")
}

func TestTeamSkills(t *testing.T) {
	t.Parallel()
	team := engine.NewTeam()
	skills := team.Skills()
	for _, skill := range skills {
		assert.Equal(t, skill.String(), "0")
	}
	team.Players[2].SetSkills(*bc.Contracts, SkillsFromString(t, "4544"))
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

func TestTrainingPointsAssigned(t *testing.T) {
	t.Parallel()
	var teamID string
	teamID = "12"
	// teamID := big.NewInt(12)
	trainingPerFieldPos := storage.TrainingPerFieldPos{10, 10, 10, 10, 10}
	trainingPerFieldPosSpecialPlayer := storage.TrainingPerFieldPos{11, 11, 11, 11, 11}
	training := storage.Training{
		teamID,
		3,
		trainingPerFieldPos,
		trainingPerFieldPos,
		trainingPerFieldPos,
		trainingPerFieldPos,
		trainingPerFieldPosSpecialPlayer,
	}
	// the sum of assigned TPs is 50, so if we had 50 available => expect no errors:
	errs := engine.CheckTraining(50, training)
	assert.Equal(t, engine.IsTrainingCorrect(50, training), true)
	assert.Equal(t, errs[0], false) // errTooMany
	assert.Equal(t, errs[1], false) // errTooManyOneSkill
	assert.Equal(t, errs[2], false) // errSpecialPlayer
	// if we only had 49 available... failure in the sums, but not in perSkill
	// the special player had (49*11)/10 = 53.9 = 53 available, so he should fail too
	errs = engine.CheckTraining(49, training)
	assert.Equal(t, engine.IsTrainingCorrect(49, training), false)
	assert.Equal(t, errs[0], true)  // errTooMany
	assert.Equal(t, errs[1], false) // errTooManyOneSkill
	assert.Equal(t, errs[2], true)  // errSpecialPlayer
	// check per skill now
	trainingPerFieldPosWrong := storage.TrainingPerFieldPos{42, 2, 2, 2, 2}
	training = storage.Training{
		teamID,
		3,
		trainingPerFieldPosWrong,
		trainingPerFieldPos,
		trainingPerFieldPos,
		trainingPerFieldPos,
		trainingPerFieldPosSpecialPlayer,
	}
	errs = engine.CheckTraining(50, training)
	assert.Equal(t, engine.IsTrainingCorrect(50, training), false)
	assert.Equal(t, errs[0], false) // errTooMany
	assert.Equal(t, errs[1], true)  // errTooManyOneSkill
	assert.Equal(t, errs[2], false) // errSpecialPlayer}
}
