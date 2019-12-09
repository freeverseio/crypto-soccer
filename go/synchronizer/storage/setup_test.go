package storage_test

import (
	"log"
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

var s *storage.Storage

func TestMain(m *testing.M) {
	var err error
	s, err = storage.NewPostgres("postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}
