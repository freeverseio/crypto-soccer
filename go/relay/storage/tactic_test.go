package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
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

func TestTacticsByVerse(t *testing.T) {
	if err := db.Begin(); err != nil {
		t.Fatal(err)
	}
	defer db.Rollback()
	tactics, err := db.TacticsByVerse(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(tactics) != 0 {
		t.Fatalf("Tactics of verse 0 are %v", len(tactics))
	}
	tactic0 := storage.Tactic{}
	tactic0.TeamID = "1"
	tactic0.ExtraAttack1 = true
	if err = db.TacticCreate(&tactic0); err != nil {
		t.Fatal(err)
	}
	tactic1 := storage.Tactic{}
	tactic1.TeamID = "2"
	tactic1.ExtraAttack2 = true
	if err = db.TacticCreate(&tactic1); err != nil {
		t.Fatal(err)
	}
	tactics, err = db.TacticsByVerse(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(tactics) != 0 {
		t.Fatalf("Tactics of verse 0 are %v", len(tactics))
	}
	if err = db.CloseVerse(); err != nil {
		t.Fatal(err)
	}
	tactics, err = db.TacticsByVerse(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(tactics) != 2 {
		t.Fatalf("Tactics of verse are %v", len(tactics))
	}
}

func TestGetRawsTactics(t *testing.T) {
	err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Rollback()

}
