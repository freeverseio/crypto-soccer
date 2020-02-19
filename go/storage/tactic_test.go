package storage_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/storage"
	"gotest.tools/assert"
)

func TestTacticNoSubstitution(t *testing.T) {
	value, err := bc.Contracts.Engine.NOSUBST(&bind.CallOpts{})
	assert.NilError(t, err)
	assert.Equal(t, value, uint8(storage.NoSubstitution))
}
