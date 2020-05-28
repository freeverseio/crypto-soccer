package contracts_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/storage"
	"gotest.tools/assert"
)

func TestPostgresNewContracts(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()
	_, err = contracts.NewFromStorage(bc.Client, tx)
	assert.Error(t, err, "no proxy address in the storage")

	assert.NilError(t, (storage.Param{contracts.ProxyName, "0x0"}).InsertOrUpdate(tx))

	_, err = contracts.NewFromStorage(bc.Client, tx)
	assert.Error(t, err, "no contract code at given address")
}

func TestPostgresNewContractsToStorage(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()
	assert.Assert(t, bc.ProxyAddress != "")
	assert.NilError(t, bc.ToStorage(tx))

	param, err := storage.ParamByName(tx, contracts.ProxyName)
	assert.NilError(t, err)
	assert.Equal(t, param.Name, contracts.ProxyName)
	assert.Equal(t, param.Value, bc.ProxyAddress)

}
