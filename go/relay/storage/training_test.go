package storage_test

import (
	"math/big"
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
	training.TeamID = big.NewInt(4)
	err = db.CreateTraining(training)
	if err != nil {
		t.Fatal(err)
	}
}
