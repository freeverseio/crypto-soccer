package playstore_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestAssignedStateProcess(t *testing.T) {
	order := storage.NewPlaystoreOrder()
	order.State = storage.PlaystoreOrderAssetAssigned
	iapTestOn := true
	_, err := playstore.New(
		[]byte{},
		*order,
		*bc.Contracts,
		bc.Owner,
		iapTestOn,
	)
	assert.Error(t, err, "unexpected end of JSON input")
}
