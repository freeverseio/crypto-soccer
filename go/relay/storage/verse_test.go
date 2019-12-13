package storage_test

import (
	"testing"
)

func TestGetCurrentVerse(t *testing.T) {
	err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Rollback()
	verse, err := db.GetLastVerse()
	if err != nil {
		t.Fatal(err)
	}
	if verse.ID != 0 {
		t.Fatalf("Expected verse 0 received %v", verse.ID)
	}
}

func TestIncreamentVerse(t *testing.T) {
	err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Rollback()
	err = db.CloseVerse()
	if err != nil {
		t.Fatal(err)
	}
	verse1, err := db.GetLastVerse()
	if err != nil {
		t.Fatal(err)
	}
	if verse1.ID != 1 {
		t.Fatalf("Expected verse 1 received %v", verse1.ID)
	}
	err = db.CloseVerse()
	if err != nil {
		t.Fatal(err)
	}
	verse2, err := db.GetLastVerse()
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
	err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Rollback()
	verse, err := db.GetVerse(0)
	if err != nil {
		t.Fatal(err)
	}
	if verse == nil {
		t.Fatal("Expected verse 0 exists")
	}
	verse, err = db.GetVerse(1)
	if err == nil {
		t.Fatal("No error on unexistent verse")
	}
	if verse != nil {
		t.Fatalf("Expected nil received %v", verse)
	}
	if err = db.CloseVerse(); err != nil {
		t.Fatal(err)
	}
	verse, err = db.GetVerse(1)
	if err != nil {
		t.Fatal(err)
	}
	if verse == nil {
		t.Fatal("Expected verse 1 exists")
	}
}
