package storagefacade_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

var s *sql.DB
var bc *testutils.BlockchainNode

func TestMain(m *testing.M) {
	var err error
	if s, err = storage.New("postgres://freeverse:freeverse@localhost:15432/cryptosoccer?sslmode=disable"); err != nil {
		log.Fatal(err)
	}
	if bc, err = testutils.NewBlockchainNode(); err != nil {
		log.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)
	os.Exit(m.Run())
}
