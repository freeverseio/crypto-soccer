package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
)

func TestGetCurrentVerse(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	verse, err := storage.GetLastVerse(tx)
	if err != nil {
		t.Fatal(err)
	}
	if verse.ID != 0 {
		t.Fatalf("Expected verse 0 received %v", verse.ID)
	}
}

func TestIncreamentVerse(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	err = storage.CloseVerse(tx)
	if err != nil {
		t.Fatal(err)
	}
	verse1, err := storage.GetLastVerse(tx)
	if err != nil {
		t.Fatal(err)
	}
	if verse1.ID != 1 {
		t.Fatalf("Expected verse 1 received %v", verse1.ID)
	}
	err = storage.CloseVerse(tx)
	if err != nil {
		t.Fatal(err)
	}
	verse2, err := storage.GetLastVerse(tx)
	if err != nil {
		t.Fatal(err)
	}
	if verse2.ID != 2 {
		t.Fatalf("Expected verse 2 received %v", verse2.ID)
	}
	if verse1.StartAt.After(verse2.StartAt) {
		t.Fatal("Verse 1 is after Verse 2")
	}
}

func TestGetVerse(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	verse, err := storage.GetVerse(tx, 0)
	if err != nil {
		t.Fatal(err)
	}
	if verse == nil {
		t.Fatal("Expected verse 0 exists")
	}
	verse, err = storage.GetVerse(tx, 1)
	if err == nil {
		t.Fatal("No error on unexistent verse")
	}
	if verse != nil {
		t.Fatalf("Expected nil received %v", verse)
	}
	if err = storage.CloseVerse(tx); err != nil {
		t.Fatal(err)
	}
	verse, err = storage.GetVerse(tx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if verse == nil {
		t.Fatal("Expected verse 1 exists")
	}
}
