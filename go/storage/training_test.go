package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"gotest.tools/assert"
)

func TestResetTrainings(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	assert.NilError(t, storage.ResetTrainingsByTimezone(tx, 0))
}

func TestTrainingCreate(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	createMinimumUniverse(t, tx)

	training := storage.Training{}
	training.TeamID = teamID
	err = training.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCurrentTraining(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	createMinimumUniverse(t, tx)

	training := storage.Training{}
	training.Verse = storage.UpcomingVerse
	training.TeamID = teamID
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
