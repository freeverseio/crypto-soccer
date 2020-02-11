package engine_test

import (
	"database/sql"
	"math/big"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"

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
	dump = spew.ConfigState{DisablePointerAddresses: true}
	os.Exit(m.Run())
}

func SkillsFromString(skills string) *big.Int {
	result, _ := new(big.Int).SetString(skills, 10)
	return result
}
