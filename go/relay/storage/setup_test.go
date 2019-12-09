package storage_test

import (
	"log"
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
)

var db *storage.Storage

func TestMain(m *testing.M) {
	var err error
	db, err = storage.NewSqlite3("../../../relay.db/00_schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}
