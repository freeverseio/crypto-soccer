package engine_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"

	log "github.com/sirupsen/logrus"
)

var s *sql.DB
var bc *testutils.BlockchainNode

func TestMain(m *testing.M) {
	var err error
	s, err = storage.New("postgres://freeverse:freeverse@localhost:15432/cryptosoccer?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	bc, err = testutils.NewBlockchainNode()
	if err != nil {
		log.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)
	os.Exit(m.Run())
}
