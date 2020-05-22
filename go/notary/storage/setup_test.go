package storage_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

var db *sql.DB
var dump spew.ConfigState

func TestMain(m *testing.M) {
	var err error
	db, err = storage.New("postgres://freeverse:freeverse@localhost:5432/market?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	dump = spew.ConfigState{DisablePointerAddresses: true, Indent: "\t"}

	os.Exit(m.Run())
}
