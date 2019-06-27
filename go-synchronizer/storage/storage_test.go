package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage/memory"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage/sqlite3"

	_ "github.com/lib/pq"
)

func suite(t *testing.T, storage storage.Storage) {
	t.Run("GetTeam", func(t *testing.T) {
		err := storage.TeamAdd(1, "ciao")
		if err != nil {
			t.Fatal(err)
		}
		team, err := storage.GetTeam(1)
		if err != nil {
			t.Fatal(err)
		}
		if team.Id != 1 {
			t.Fatalf("Expected 1 result %v", team.Id)
		}
		if team.Name != "ciao" {
			t.Fatalf("Expected ciao result %v", team.Name)
		}
	})
	t.Run("GetUnexistentTeam", func(t *testing.T) {
		_, err := storage.GetTeam(0)
		if err == nil {
			t.Fatal("No error on get unexistent team")
		}
	})
	t.Run("BlockNumber", func(t *testing.T) {
		blockNumber, err := storage.GetBlockNumber()
		if err != nil {
			t.Fatal(err)
		}
		if blockNumber == nil {
			t.Fatalf("Expected not nil")
		}
		if blockNumber.Int64() != -1 {
			t.Fatalf("Expected -1 result %v", blockNumber)
		}

		err = storage.SetBlockNumber(big.NewInt(3))
		if err != nil {
			t.Fatal(err)
		}

		blockNumber, err = storage.GetBlockNumber()
		if err != nil {
			t.Fatal(err)
		}
		if blockNumber.String() != "3" {
			t.Fatalf("Expected 3 result %v", blockNumber)
		}
	})
}

func TestMemory(t *testing.T) {
	storage := memory.New()
	suite(t, storage)
}

func TestSqlite3(t *testing.T) {
	storage, err := sqlite3.New("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	suite(t, storage)
}

// func TestPostgres(t *testing.T) {
// 	storage, err := postgres.New("postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	suite(t, storage)
// }
