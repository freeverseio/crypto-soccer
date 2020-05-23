package consumer_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	"github.com/freeverseio/crypto-soccer/go/testutils"

	log "github.com/sirupsen/logrus"
)

var bc *testutils.BlockchainNode
var db *sql.DB
var googleCredentials []byte

func TestMain(m *testing.M) {
	var err error
	db, err = postgres.New("postgres://freeverse:freeverse@localhost:5432/market?sslmode=disable")
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
