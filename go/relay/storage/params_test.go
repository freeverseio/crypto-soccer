package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
)

func TestGetParam(t *testing.T) {
	db, err := storage.NewSqlite3("../../../relay.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	verse, err := db.GetVerse()
	if err != nil {
		t.Fatal(err)
	}
	if verse.Uint64() != 0 {
		t.Fatalf("Wrong verse %v", verse)
	}
}

func TestSetVerse(t *testing.T) {
	db, err := storage.NewSqlite3("../../../relay.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	err = db.SetVerse(big.NewInt(1))
	if err != nil {
		t.Fatal(err)
	}
	verse, err := db.GetVerse()
	if err != nil {
		t.Fatal(err)
	}
	if verse.Uint64() != 1 {
		t.Fatalf("Wrong verse %v", verse)
	}
}
