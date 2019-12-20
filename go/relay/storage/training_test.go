package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
)

func TestTrainingCreate(t *testing.T) {
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
