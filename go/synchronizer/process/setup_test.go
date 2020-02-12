package process_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

var universedb *sql.DB
var bc *testutils.BlockchainNode
var dump spew.ConfigState
var namesdb *names.Generator

const ipfsURL = "localhost:5001"

func TestMain(m *testing.M) {
	var err error
	namesdb, err = names.New("../../names/sql/names.db")
	if err != nil {
		log.Fatal(err)
	}
	universedb, err = storage.New("postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	bc, err = testutils.NewBlockchainNode()
	if err != nil {
		log.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)
	bc.InitOneTimezone(1)
	dump = spew.ConfigState{DisablePointerAddresses: true}
	os.Exit(m.Run())
}
