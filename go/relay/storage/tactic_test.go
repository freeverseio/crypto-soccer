package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
)

func TestTacticCreate(t *testing.T) {
	err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Rollback()
	tacticID := uint8(16)
	teamId := big.NewInt(1)
	shirts := [14]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	extraAttack := [10]bool{false, false, false, false, false, false, false, false, false, false}
	substitutions := [3]uint8{11, 11, 11}
	subsRounds := [3]uint8{2, 3, 4}
	verse := uint64(10)
	err = db.TacticCreate(storage.Tactic{teamId, tacticID, shirts, extraAttack, substitutions, subsRounds}, verse)
	if err != nil {
		t.Fatal(err)
	}
	count, err := db.TacticCount(nil)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expecting 1 tactic, got %v", count)
	}

	count, err = db.TacticCount(&verse)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expecting 1 tactic, got %v", count)
	}
	tc, err := db.GetTactic(teamId, verse)
	if err != nil {
		t.Fatal(err)
	}
	if tc.TacticID != tacticID {
		t.Fatalf("expecting tacticID 1, got %v", tc.TacticID)
	}

	tc, err = db.GetTactic(big.NewInt(2), verse)
	if err == nil {
		t.Fatal("team 2 does not exist and should fail")
	}
}
