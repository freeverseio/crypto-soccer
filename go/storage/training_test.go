package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"gotest.tools/assert"
)

func TestResetTrainings(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	assert.NilError(t, storage.ResetTrainingsByTimezone(tx, 0))
}

func TestTrainingCreate(t *testing.T) {
	t.Skip("******************** REACTIVE  **********************")
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	training := storage.Training{}
	training.TeamID = "4"
	err = training.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCurrentTraining(t *testing.T) {
	t.Skip("******************** REACTIVE  **********************")
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	training := storage.Training{}
	training.Verse = storage.UpcomingVerse
	training.TeamID = "4"
	err = training.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}

	trainings, err := storage.UpcomingTrainings(tx)
	if err != nil {
		t.Fatal(err)
	}
	if len(trainings) != 1 {
		t.Fatalf("Expected 1 got %v", len(trainings))
	}

	training.Verse = 1
	err = training.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}

	trainings, err = storage.UpcomingTrainings(tx)
	if err != nil {
		t.Fatal(err)
	}
	if len(trainings) != 1 {
		t.Fatalf("Expected 1 got %v", len(trainings))
	}

}
