package playstore_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestMachineCreation(t *testing.T) {
	client := NewMockClientService()
	iapTestOn := true

	order := storage.NewPlaystoreOrder()
	m, err := playstore.New(
		client,
		*order,
		*bc.Contracts,
		bc.Owner,
		iapTestOn,
	)
	assert.NilError(t, err)
	assert.NilError(t, m.Process())
	assert.Equal(t, m.Order().State, storage.PlaystoreOrderFailed)
}
