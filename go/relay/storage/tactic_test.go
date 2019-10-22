package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
)

func TestTacticCreate(t *testing.T) {
	db, err := storage.NewSqlite3("../../../relay.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	teamId := big.NewInt(1)
	shirts := [11]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	extraAttack := [10]bool{false, false, false, false, false, false, false, false, false, false}
	err = db.TacticCreate(storage.Tactic{teamId, 4, 3, 3, shirts, extraAttack})
	if err != nil {
		t.Fatal(err)
	}
	count, err := db.TacticCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expecting 1 tactic, got %v", count)
	}
	tc, err := db.GetTactic(teamId)
	if err != nil {
		t.Fatal(err)
	}
	if tc.Defense != 4 || tc.Center != 3 || tc.Attack != 3 {
		t.Fatalf("expecting 4-3-3, got %v-%v-%v", tc.Defense, tc.Center, tc.Attack)
	}

	tc, err = db.GetTactic(big.NewInt(2))
	if err == nil {
		t.Fatal("team 2 does not exist and should fail")
	}
}
