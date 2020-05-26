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
