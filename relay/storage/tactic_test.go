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
	shirts := [11]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	extraAttack := [10]uint8{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	err = db.TacticCreate(storage.Tactic{teamId, 4, 3, 3, shirts, extraAttack})
	if err != nil {
		t.Fatal(err)
	}
}
