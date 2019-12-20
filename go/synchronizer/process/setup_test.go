package process_test

import (
	"log"
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

var universedb *storage.Storage

func TestMain(m *testing.M) {
	var err error
	universedb, err = storage.NewPostgres("postgres://freeverse:freeverse@localhost:15432/cryptosoccer?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
