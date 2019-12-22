package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
)

func TestTacticByTeamID(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	tc, err := storage.TacticByTeamID(tx, "5")
	if err != nil {
		t.Fatal(err)
	}
	if tc != nil {
		t.Fatalf("Received %v", tc)
	}
}

func TestTacticCreate(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	count, err := storage.TacticCount(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expecting 0 tactic, got %v", count)
	}
	tactic := storage.DefaultTactic("16")
	if tactic.TeamID != "16" {
		t.Fatalf("Expected 16 but %v", tactic.TeamID)
	}
	err = tactic.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}
	count, err = storage.TacticCount(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expecting 1 tactic, got %v", count)
	}
	t.Logf("TacticId: %v", tactic.TeamID)
	tc, err := storage.TacticByTeamID(tx, tactic.TeamID)
	if err != nil {
		t.Fatal(err)
	}
	if *tc != *tactic {
		t.Fatalf("expecting tacticID %v, got %v", tactic, tc)
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
	tactic0.Verse = 1
	tactic0.TeamID = "1"
	tactic0.ExtraAttack1 = true
	if err = tactic0.Insert(tx); err != nil {
		t.Fatal(err)
	}
	tactic1 := storage.Tactic{}
	tactic1.Verse = 1
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
		t.Fatalf("Tactics of verse 1 are %v", len(tactics))
	}
}
