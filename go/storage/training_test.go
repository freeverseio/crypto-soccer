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

func TestTrainingReset(t *testing.T) {
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	createMinimumUniverse(t, tx)

	training := storage.NewTraining()
	training.TeamID = teamID
	training.Defenders.Speed = 3
	assert.NilError(t, training.Insert(tx))

	trainings, err := storage.TrainingsByTimezone(tx, int(timezoneIdx))
	assert.NilError(t, err)
	assert.Equal(t, len(trainings), 1)
	assert.Equal(t, training.Defenders.Speed, trainings[0].Defenders.Speed)

	assert.NilError(t, storage.ResetTrainingsByTimezone(tx, timezoneIdx))
	trainings, err = storage.TrainingsByTimezone(tx, int(timezoneIdx))
	assert.NilError(t, err)
	assert.Equal(t, len(trainings), 1)
	assert.Equal(t, 0, trainings[0].Defenders.Speed)
}

func TestTrainingCreteDefaultPerTimezone(t *testing.T) {
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	createMinimumUniverse(t, tx)
	trainings, err := storage.TrainingsByTimezone(tx, int(timezoneIdx))
	assert.NilError(t, err)
	assert.Equal(t, len(trainings), 0)

	assert.NilError(t, storage.CreateDefaultTrainingByTimezone(tx, timezoneIdx+1))
	trainings, err = storage.TrainingsByTimezone(tx, int(timezoneIdx+1))
	assert.NilError(t, err)
	assert.Equal(t, len(trainings), 0)

	assert.NilError(t, storage.CreateDefaultTrainingByTimezone(tx, timezoneIdx))
	trainings, err = storage.TrainingsByTimezone(tx, int(timezoneIdx))
	assert.NilError(t, err)
	assert.Equal(t, len(trainings), 2)
}
