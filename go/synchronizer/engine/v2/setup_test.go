package v2_test

import (
	"database/sql"
	"math/big"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
	"gotest.tools/assert"

	log "github.com/sirupsen/logrus"
)

var db *sql.DB
var bc *testutils.BlockchainNode
var dump spew.ConfigState

func TestMain(m *testing.M) {
	var err error
	db, err = storage.New("postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	bc, err = testutils.NewBlockchain()
	if err != nil {
		log.Fatal(err)
	}
	dump = spew.ConfigState{DisablePointerAddresses: true, Indent: "\t"}
	os.Exit(m.Run())
}

func SkillsFromString(t *testing.T, skills string) *big.Int {
	result, ok := new(big.Int).SetString(skills, 10)
	assert.Equal(t, ok, true)
	return result
}

func TestSkillsFromString(t *testing.T) {
	skills := SkillsFromString(t, "40439920000726868070503716865792521545121682176182486071370780491777")
	assert.Equal(t, skills.String(), "40439920000726868070503716865792521545121682176182486071370780491777")
}

func IsTrainingCorrect(availableTPs int, tr storage.Training) bool {
	err := CheckTraining(availableTPs, tr)
	return !(err[0] || err[1] || err[2])
}

func CheckTraining(availableTPs int, tr storage.Training) [3]bool {
	errTooMany := false
	errTooManyOneSkill := false
	err := CheckTrainingPerFieldPos(availableTPs, tr.Goalkeepers)
	errTooMany = errTooMany || err[0]
	errTooManyOneSkill = errTooManyOneSkill || err[1]
	err = CheckTrainingPerFieldPos(availableTPs, tr.Defenders)
	errTooMany = errTooMany || err[0]
	errTooManyOneSkill = errTooManyOneSkill || err[1]
	err = CheckTrainingPerFieldPos(availableTPs, tr.Midfielders)
	errTooMany = errTooMany || err[0]
	errTooManyOneSkill = errTooManyOneSkill || err[1]
	err = CheckTrainingPerFieldPos(availableTPs, tr.Attackers)
	errTooMany = errTooMany || err[0]
	errTooManyOneSkill = errTooManyOneSkill || err[1]
	// Special Player has extra 10% points, calculated in this precise integer-division manner:
	availableTPs = (availableTPs * 11) / 10
	err = CheckTrainingPerFieldPos(availableTPs, tr.SpecialPlayer)
	errSpecialPlayer := err[0] || err[1]
	return [3]bool{errTooMany, errTooManyOneSkill, errSpecialPlayer}
}

func CheckTrainingPerFieldPos(availableTPs int, tr storage.TrainingPerFieldPos) [2]bool {
	sum := tr.Shoot + tr.Speed + tr.Pass + tr.Defence + tr.Endurance
	errTooMany := (sum > availableTPs)
	errTooManyOneSkill := (100*tr.Shoot > 60*availableTPs) ||
		(100*tr.Speed > 60*availableTPs) ||
		(100*tr.Pass > 60*availableTPs) ||
		(100*tr.Defence > 60*availableTPs) ||
		(100*tr.Endurance > 60*availableTPs)
	return [2]bool{errTooMany, errTooManyOneSkill}
}
