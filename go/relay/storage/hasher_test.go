package storage_test

import (
	"encoding/hex"
	"testing"
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

func TestHashVerseOfEmptyVerse(t *testing.T) {
	err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Rollback()
	err = db.CloseVerse()
	if err != nil {
		t.Fatal(err)
	}
	hash, err := db.HashVerse(1)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(hash) != "2dba5dbc339e7316aea2683faf839c1b7b1ee2313db792112588118df066aa35" {
		t.Fatalf("Wrong result %v", hex.EncodeToString(hash))
	}
	err = db.CloseVerse()
	if err != nil {
		t.Fatal(err)
	}
	hash, err = db.HashVerse(2)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(hash) != "2dba5dbc339e7316aea2683faf839c1b7b1ee2313db792112588118df066aa35" {
		t.Fatalf("Wrong result %v", hex.EncodeToString(hash))
	}
}

func TestHashVerse(t *testing.T) {
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
	if err = db.CloseVerse(); err != nil {
		t.Fatal(err)
	}
	hash, err := db.HashVerse(1)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(hash) == "5df6e0e2761359d30a8275058e299fcc0381534545f55cf43e41983f5d4c9456" {
		t.Fatal("Empty verse hash")
	}
}
