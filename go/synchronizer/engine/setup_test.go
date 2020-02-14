package engine_test

import (
	"database/sql"
	"math/big"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/freeverseio/crypto-soccer/go/relay/storage"
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
	bc, err = testutils.NewBlockchainNode()
	if err != nil {
		log.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)
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
