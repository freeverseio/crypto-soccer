package relay_test

import (
	//"math/big"
	"testing"

	//"github.com/ethereum/go-ethereum/accounts/abi/bind"
	//"github.com/ethereum/go-ethereum/core/types"
	//"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/relay"
	"gotest.tools/assert"
)

func TestSubmitActionRoot(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auth := bind.NewKeyedTransactor(bc.Owner)
	p, err := relay.NewProcessor(bc.Client, auth, db, bc.Contracts.Updates, "localhost:5001")
	assert.NilError(t, err)
	assert.NilError(t, p.Process(tx))
}
