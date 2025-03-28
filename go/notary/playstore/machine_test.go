package playstore_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/googleplaystoreutils"
	"github.com/freeverseio/crypto-soccer/go/notary/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestMachineCreation(t *testing.T) {
	client := googleplaystoreutils.NewMockClientService()
	iapTestOn := true
	order := storage.NewPlaystoreOrder()
	_, err := playstore.New(
		client,
		*order,
		*bc.Contracts,
		bc.Owner,
		namesdb,
		iapTestOn,
	)
	assert.NilError(t, err)
}

func TestMachineCreationFailedState(t *testing.T) {
	client := googleplaystoreutils.NewMockClientService()
	iapTestOn := true
	order := storage.NewPlaystoreOrder()
	order.State = storage.PlaystoreOrderFailed
	m, err := playstore.New(
		client,
		*order,
		*bc.Contracts,
		bc.Owner,
		namesdb,
		iapTestOn,
	)
	assert.NilError(t, err)
	assert.NilError(t, m.Process())
}

func TestMachineCreationRefundedState(t *testing.T) {
	client := googleplaystoreutils.NewMockClientService()
	iapTestOn := true
	order := storage.NewPlaystoreOrder()
	order.State = storage.PlaystoreOrderRefunded
	m, err := playstore.New(
		client,
		*order,
		*bc.Contracts,
		bc.Owner,
		namesdb,
		iapTestOn,
	)
	assert.NilError(t, err)
	assert.NilError(t, m.Process())
}

func TestMachineCreationCompleteState(t *testing.T) {
	client := googleplaystoreutils.NewMockClientService()
	iapTestOn := true
	order := storage.NewPlaystoreOrder()
	order.State = storage.PlaystoreOrderComplete
	m, err := playstore.New(
		client,
		*order,
		*bc.Contracts,
		bc.Owner,
		namesdb,
		iapTestOn,
	)
	assert.NilError(t, err)
	assert.NilError(t, m.Process())
}
