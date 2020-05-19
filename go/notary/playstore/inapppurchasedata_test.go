package playstore_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/playstore"
	"gotest.tools/assert"
)

func TestInappPurchaseDataFromReceipt(t *testing.T) {
	receipt := ""
	data, err := playstore.InappPurchaseDataFromReceipt(receipt)
	assert.NilError(t, err)
	assert.Equal(t, data.OrderId, "")
}
