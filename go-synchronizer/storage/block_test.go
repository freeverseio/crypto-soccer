package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
)

func TestGetBlockNumber(t *testing.T) {
	storage, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	number, err := storage.GetBlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	if number != 0 {
		t.Fatalf("Expected 0 result %v", number)
	}
}

func TestSetBlockNumber(t *testing.T) {
	storage, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	err = storage.SetBlockNumber(3)
	if err != nil {
		t.Fatal(err)
	}
	number, err := storage.GetBlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	if number != 3 {
		t.Fatalf("Expected 3 result %v", number)
	}
}
