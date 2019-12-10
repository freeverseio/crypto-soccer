package storage_test

import "testing"

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
	verse, err := db.GetLastVerse()
	if err != nil {
		t.Fatal(err)
	}
	if verse.ID != 1 {
		t.Fatalf("Expected verse 1 received %v", verse.ID)
	}
}
