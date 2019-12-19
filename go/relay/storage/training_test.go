package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
)

func TestTrainingCreate(t *testing.T) {
	err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Rollback()
	training := storage.Training{}
	training.TeamID = "4"
	err = db.CreateTraining(training)
	if err != nil {
		t.Fatal(err)
	}
}
