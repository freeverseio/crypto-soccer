package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
)

func TestTrainingCreate(t *testing.T) {
	db, err := storage.NewSqlite3("../../../relay.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	training := storage.Training{}
	training.TeamID = big.NewInt(4)
	err = db.CreateTraining(training)
	if err != nil {
		t.Fatal(err)
	}
}
