package consumer_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/relay/consumer"
	"gotest.tools/assert"
)

func TestSubmitActionRoot(t *testing.T) {
	t.Parallel()
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auth := bind.NewKeyedTransactor(bc.Owner)
	p := consumer.NewProcessor(bc.Client, auth, bc.Contracts.Updates, "localhost:5001")
	assert.NilError(t, p.Process(tx))
}
