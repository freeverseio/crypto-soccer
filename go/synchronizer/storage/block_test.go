package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

func TestGetBlockNumber(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	number, err := storage.GetBlockNumber(tx)
	if err != nil {
		t.Fatal(err)
	}
	if number != 0 {
		t.Fatalf("Expected 0 result %v", number)
	}
}

func TestSetBlockNumber(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	err = storage.SetBlockNumber(tx, 3)
	if err != nil {
		t.Fatal(err)
	}
	number, err := storage.GetBlockNumber(tx)
	if err != nil {
		t.Fatal(err)
	}
	if number != 3 {
		t.Fatalf("Expected 3 result %v", number)
	}
}
