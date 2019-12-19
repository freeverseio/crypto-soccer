package relay_test

import (
	"log"
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
)

var db *storage.Storage

func TestMain(m *testing.M) {
	var err error
	db, err = storage.NewPostgres("postgres://freeverse:freeverse@localhost:15433/relay?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}
