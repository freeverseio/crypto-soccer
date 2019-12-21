package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
)

func TestVerseTacticCreate(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	tactic := storage.VerseTactic{}
	tactic.Verse = 1
	tactic.Tactic = *storage.DefaultTactic("16")
	err = tactic.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}
	count, err := storage.VerseTacticCount(tx, 0)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expecting 0 tactic, got %v", count)
	}

	count, err = storage.VerseTacticCount(tx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expecting 1 tactic, got %v", count)
	}
	tc, err := storage.VerseTacticByTeamIDAndVerse(tx, tactic.Tactic.TeamID, tactic.Verse)
	if err != nil {
		t.Fatal(err)
	}
	if *tc != tactic {
		t.Fatalf("expecting tacticID %v, got %v", tactic, tc)
	}

	tc, err = storage.VerseTacticByTeamIDAndVerse(tx, "2", tactic.Verse)
	if err == nil {
		t.Fatal("team 2 does not exist and should fail")
	}
}

func TestVerseTacticsByVerse(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	tactics, err := storage.VerseTacticsByVerse(tx, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(tactics) != 0 {
		t.Fatalf("Tactics of verse 0 are %v", len(tactics))
	}
	tactic0 := storage.VerseTactic{}
	tactic0.Verse = 1
	tactic0.Tactic.TeamID = "1"
	tactic0.Tactic.ExtraAttack1 = true
	if err = tactic0.Insert(tx); err != nil {
		t.Fatal(err)
	}
	tactic1 := storage.VerseTactic{}
	tactic1.Verse = 1
	tactic1.Tactic.TeamID = "2"
	tactic1.Tactic.ExtraAttack2 = true
	if err = tactic1.Insert(tx); err != nil {
		t.Fatal(err)
	}
	tactics, err = storage.VerseTacticsByVerse(tx, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(tactics) != 0 {
		t.Fatalf("Tactics of verse 0 are %v", len(tactics))
	}
	if err = storage.CloseVerse(tx); err != nil {
		t.Fatal(err)
	}
	tactics, err = storage.VerseTacticsByVerse(tx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(tactics) != 2 {
		t.Fatalf("Tactics of verse 1 are %v", len(tactics))
	}
}
