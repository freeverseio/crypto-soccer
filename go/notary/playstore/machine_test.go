package playstore_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestMachineCreation(t *testing.T) {
	order := storage.NewPlaystoreOrder()
	machine := playstore.New([]byte{}, *order)
	assert.NilError(t, machine.Process())
	assert.Equal(t, machine.Order().State, storage.PlaystoreOrderFailed)
}
