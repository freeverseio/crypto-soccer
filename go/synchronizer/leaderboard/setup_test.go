package leaderboard_test

import (
	"database/sql"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/useractions"
	"github.com/freeverseio/crypto-soccer/go/useractions/ipfs"

	"github.com/davecgh/go-spew/spew"
	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
	"gotest.tools/assert"
)

var universedb *sql.DB
var bc *testutils.BlockchainNode
var dump spew.ConfigState
var namesdb *names.Generator
var useractionsPublishService useractions.UserActionsPublishService

func TestMain(m *testing.M) {
	var err error
	namesdb, err = names.New("../../names/sql/names.db")
	if err != nil {
		log.Fatal(err)
	}
	universedb, err = storage.New("postgres://freeverse:freeverse@crypto-soccer_devcontainer_dockerhost_1:5432/cryptosoccer?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	bc, err = testutils.NewBlockchain()
	if err != nil {
		log.Fatal(err)
	}
	useractionsPublishService = ipfs.NewUserActionsPublishService("/ip4/127.0.0.1/tcp/5001")
	dump = spew.ConfigState{DisablePointerAddresses: true}
	os.Exit(m.Run())
}

func SkillsFromString(t *testing.T, skills string) *big.Int {
	result, ok := new(big.Int).SetString(skills, 10)
	assert.Equal(t, ok, true)
	return result
}
