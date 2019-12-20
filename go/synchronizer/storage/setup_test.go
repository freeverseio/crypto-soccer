package storage_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

var s *sql.DB

func TestMain(m *testing.M) {
	var err error
	s, err = storage.New("postgres://freeverse:freeverse@localhost:15432/cryptosoccer?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}
