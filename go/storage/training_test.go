package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"gotest.tools/assert"
)

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

	trainings, err := storage.UpcomingTrainings(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(trainings), 0)

	training := storage.Training{}
	training.Verse = storage.UpcomingVerse
	training.TeamID = teamID
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

func TestTrainingResetTrainings(t *testing.T) {
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	createMinimumUniverse(t, tx)

	training := storage.NewTraining()
	training.Timezone = int(timezoneIdx)
	training.TeamID = teamID
	training.DefendersDefence = 5
	assert.NilError(t, training.Insert(tx))
	trainings, err := storage.UpcomingTrainings(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(trainings), 1)
	assert.Equal(t, trainings[0].DefendersDefence, 5)
	assert.NilError(t, storage.ResetTrainingsByTimezone(tx, timezoneIdx))
	trainings, err = storage.UpcomingTrainings(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(trainings), 1)
	assert.Equal(t, trainings[0].DefendersDefence, 0)
}
