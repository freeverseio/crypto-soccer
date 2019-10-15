package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/relay/storage"
)

func TestTacticCreate(t *testing.T) {
	db, err := storage.NewSqlite3("../db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	teamId := big.NewInt(1)
	err = db.TacticCreate(storage.Tactic{teamId, 4, 3, 3, nil})
	if err != nil {
		t.Fatal(err)
	}
}
