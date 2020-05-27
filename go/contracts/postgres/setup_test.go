package postgres_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

var universedb *sql.DB
var bc *testutils.BlockchainNode

const ipfsURL = "/ip4/127.0.0.1/tcp/5001"

func TestMain(m *testing.M) {
	var err error
	universedb, err = storage.New("postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	bc, err = testutils.NewBlockchain()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
