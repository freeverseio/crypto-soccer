package storage_test

import (
	"testing"
)

func TestTacticCreate(t *testing.T) {
	err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Rollback()
	tactic := db.DefaultTactic("16")
	err = db.TacticCreate(tactic)
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

	verse := uint64(3)
	count, err = db.TacticCount(&verse)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expecting 1 tactic, got %v", count)
	}
	tc, err := db.GetTactic(tactic.TeamID, verse)
	if err != nil {
		t.Fatal(err)
	}
	if uint8(tc.TacticID) != uint8(tactic.TacticID) {
		t.Fatalf("expecting tacticID 1, got %v", tc.TacticID)
	}

	tc, err = db.GetTactic("2", verse)
	if err == nil {
		t.Fatal("team 2 does not exist and should fail")
	}
}

func TestGetRawsTactics(t *testing.T) {
	err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Rollback()

}
