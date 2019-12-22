package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
)

func TestVerseTrainingInsert(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	training := storage.VerseTraining{}
	training.Verse = 1
	training.Training.TeamID = "4"
	err = training.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}

	result, err := storage.VerseTrainingByVerse(tx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("Expected 1 received %v", len(result))
	}

	training.Verse = 2
	err = training.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}

	result, err = storage.VerseTrainingByVerse(tx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("Expected 1 received %v", len(result))
	}
}
