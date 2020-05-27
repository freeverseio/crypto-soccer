package postgres_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/contracts/postgres"
	"github.com/freeverseio/crypto-soccer/go/storage"
	"gotest.tools/assert"
)

func TestPostgresNewContracts(t *testing.T) {
	tx, err := universedb.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()
	_, err = postgres.NewFromStorage(bc.Client, tx)
	assert.Error(t, err, "no contract code at given address")

	assert.NilError(t, postgres.ToStorage(tx, bc.Contracts))

	contracts, err := postgres.NewFromStorage(bc.Client, tx)
	assert.NilError(t, err)
	assert.Assert(t, contracts != nil)
	assert.Equal(t, contracts.ProxyAddress, bc.Contracts.ProxyAddress)
}

func TestPostgresNewContractsToStorage(t *testing.T) {
	tx, err := universedb.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()
	proxyAddress := bc.Contracts.ProxyAddress
	assert.Assert(t, proxyAddress != "")
	assert.NilError(t, postgres.ToStorage(tx, bc.Contracts))

	param, err := storage.ParamByName(tx, "PROXY")
	assert.NilError(t, err)
	assert.Equal(t, param.Name, "PROXY")
	assert.Equal(t, param.Value, proxyAddress)

}
