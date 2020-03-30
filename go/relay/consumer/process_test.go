package consumer_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"gotest.tools/assert"
)

func TestSubmitActionRoot(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auth := bind.NewKeyedTransactor(bc.Owner)
	p, err := relay.NewProcessor(bc.Client, auth, bc.Contracts.Updates, "localhost:5001")
	assert.NilError(t, err)
	assert.NilError(t, p.Process(tx))
}
