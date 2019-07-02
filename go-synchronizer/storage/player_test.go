package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"

	_ "github.com/lib/pq"
)

func TestPlayerCount(t *testing.T) {
	storage, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	count, err := storage.PlayerCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 result %v", count)
	}
}

func TestPlayerAdd(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	player := storage.Player{3, big.NewInt(43524)}
	err = sto.PlayerAdd(&player)
	if err != nil {
		t.Fatal(err)
	}
	count, err := sto.PlayerCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 result %v", count)
	}
}
