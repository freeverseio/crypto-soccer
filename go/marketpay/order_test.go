package marketpay_test

import (
	"encoding/json"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/marketpay"
	"gotest.tools/assert"
)

func TestOrderMarshallNumber(t *testing.T) {
	order := marketpay.Order{}
	_, err := json.Marshal(order)
	assert.NilError(t, err)

	order.Amount = "102"
	_, err = json.Marshal(order)
	assert.NilError(t, err)
}
