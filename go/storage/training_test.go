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

func TestTrainingTrainingsByTimezone(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	createMinimumUniverse(t, tx)

	trainings, err := storage.TrainingsByTimezone(tx, int(timezoneIdx))
	assert.NilError(t, err)
	assert.Equal(t, len(trainings), 0)

	training := storage.NewTraining()
	training.TeamID = teamID
	assert.NilError(t, training.Insert(tx))
	trainings, err = storage.TrainingsByTimezone(tx, int(timezoneIdx))
	assert.NilError(t, err)
	assert.Equal(t, len(trainings), 1)
}

func TestTrainingResetTrainingsInTimezone(t *testing.T) {
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	createMinimumUniverse(t, tx)

	assert.NilError(t, storage.ResetTrainingsByTimezone(tx, timezoneIdx))

	training := storage.NewTraining()
	training.TeamID = teamID
	training.DefendersPass = 4
	assert.NilError(t, training.Insert(tx))

	assert.NilError(t, storage.ResetTrainingsByTimezone(tx, timezoneIdx))
	trainings, err := storage.TrainingsByTimezone(tx, int(timezoneIdx))
	assert.NilError(t, err)
	assert.Equal(t, len(trainings), 1)
	assert.Equal(t, trainings[0].DefendersPass, 0)
}
