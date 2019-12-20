package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
)

func TestTacticCreate(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	tactic := storage.DefaultTactic("16")
	err = tactic.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}
	count, err := storage.TacticCount(tx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expecting 1 tactic, got %v", count)
	}

	verse := uint64(3)
	count, err = storage.TacticCount(tx, &verse)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expecting 1 tactic, got %v", count)
	}
	tc, err := storage.TacticByTeamIDAndVerse(tx, tactic.TeamID, verse)
	if err != nil {
		t.Fatal(err)
	}
	if uint8(tc.TacticID) != uint8(tactic.TacticID) {
		t.Fatalf("expecting tacticID 1, got %v", tc.TacticID)
	}

	tc, err = storage.TacticByTeamIDAndVerse(tx, "2", verse)
	if err == nil {
		t.Fatal("team 2 does not exist and should fail")
	}
}

func TestTacticsByVerse(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	tactics, err := storage.TacticsByVerse(tx, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(tactics) != 0 {
		t.Fatalf("Tactics of verse 0 are %v", len(tactics))
	}
	tactic0 := storage.Tactic{}
	tactic0.TeamID = "1"
	tactic0.ExtraAttack1 = true
	if err = tactic0.Insert(tx); err != nil {
		t.Fatal(err)
	}
	tactic1 := storage.Tactic{}
	tactic1.TeamID = "2"
	tactic1.ExtraAttack2 = true
	if err = tactic1.Insert(tx); err != nil {
		t.Fatal(err)
	}
	tactics, err = storage.TacticsByVerse(tx, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(tactics) != 0 {
		t.Fatalf("Tactics of verse 0 are %v", len(tactics))
	}
	if err = storage.CloseVerse(tx); err != nil {
		t.Fatal(err)
	}
	tactics, err = storage.TacticsByVerse(tx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(tactics) != 2 {
		t.Fatalf("Tactics of verse are %v", len(tactics))
	}
}
