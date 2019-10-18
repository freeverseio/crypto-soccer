package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/relay/storage"
	_ "github.com/lib/pq"
)

func TestNew(t *testing.T) {
	_, err := storage.NewSqlite3("../db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
}
