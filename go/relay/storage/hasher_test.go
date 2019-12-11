package storage_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
)

func TestHashVerseOfEmptyDB(t *testing.T) {
	err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Rollback()
	hash, err := db.HashVerse(0)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(hash) != "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		t.Fatalf("Wrong result %v", hex.EncodeToString(hash))
	}
	hash, err = db.HashVerse(1)
	if err == nil {
		t.Fatal("Expected error on hashing unexistent verse")
	}
}

func TestHashVerse(t *testing.T) {
	err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Rollback()
	tactic := storage.Tactic{}
	tactic.TeamID = big.NewInt(2)
	err = db.TacticCreate(tactic, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err = db.CloseVerse(); err != nil {
		t.Fatal(err)
	}
	hash, err := db.HashVerse(1)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(hash) != "5df6e0e2761359d30a8275058e299fcc0381534545f55cf43e41983f5d4c9456" {
		t.Fatalf("Wrong result %v", hex.EncodeToString(hash))
	}
}
