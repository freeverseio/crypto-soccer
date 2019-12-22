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
