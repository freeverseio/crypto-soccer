package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
	"gotest.tools/assert"
)

func TestTacticCreate(t *testing.T) {
	t.Skip("******************** REACTIVE  **********************")
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
	tactic := storage.DefaultTactic("16", 2)
	if tactic.TeamID != "16" {
		t.Fatalf("Expected 16 but %v", tactic.TeamID)
	}
	err = tactic.Insert(tx)
	assert.NilError(t, err)
	count, err = storage.TacticCount(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expecting 1 tactic, got %v", count)
	}
	tc, err := storage.TacticByTeamID(tx, tactic.TeamID)
	assert.NilError(t, err)
	assert.Equal(t, *tc, *tactic)
}
