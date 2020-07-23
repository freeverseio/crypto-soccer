package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"gotest.tools/assert"
)

func TestParamByName(t *testing.T) {
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	param, err := storage.ParamByName(tx, "prova")
	assert.NilError(t, err)
	assert.Assert(t, param == nil)
}

func TestParamInsertOrUpdate(t *testing.T) {
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	param := storage.Param{}
	param.Name = "prova"
	param.Value = "3"
	assert.NilError(t, param.InsertOrUpdate(tx))

	result, err := storage.ParamByName(tx, param.Name)
	assert.NilError(t, err)
	assert.Equal(t, *result, param)

	param.Value = "4"
	assert.NilError(t, param.InsertOrUpdate(tx))

	result, err = storage.ParamByName(tx, param.Name)
	assert.NilError(t, err)
	assert.Equal(t, *result, param)
}

func TestParamParams(t *testing.T) {
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	params, err := storage.Params(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(params), 0)

	param := storage.Param{}
	param.Name = "prova"
	param.Value = "3"
	assert.NilError(t, param.InsertOrUpdate(tx))
	param.Value = "4"
	assert.NilError(t, param.InsertOrUpdate(tx))

	params, err = storage.Params(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(params), 1)

	param.Name = "riprova"
	assert.NilError(t, param.InsertOrUpdate(tx))
	params, err = storage.Params(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(params), 2)
}

func TestGetBlockNumber(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	number, err := storage.GetBlockNumber(tx)
	if err != nil {
		t.Fatal(err)
	}
	if number != 0 {
		t.Fatalf("Expected 0 result %v", number)
	}
}

func TestSetBlockNumber(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	err = storage.SetBlockNumber(tx, 3)
	if err != nil {
		t.Fatal(err)
	}
	number, err := storage.GetBlockNumber(tx)
	if err != nil {
		t.Fatal(err)
	}
	if number != 3 {
		t.Fatalf("Expected 3 result %v", number)
	}
}
