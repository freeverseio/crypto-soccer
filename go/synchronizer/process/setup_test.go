package process_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

var universedb *sql.DB
var relaydb *sql.DB

func TestMain(m *testing.M) {
	var err error
	universedb, err = storage.New("postgres://freeverse:freeverse@localhost:15432/cryptosoccer?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	relaydb, err = storage.New("postgres://freeverse:freeverse@localhost:15433/relay?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
