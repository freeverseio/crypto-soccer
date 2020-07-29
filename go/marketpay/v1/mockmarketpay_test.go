package v1_test

import (
	"testing"

	v1 "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"gotest.tools/assert"
)

func TestMockMarketPaySetStatus(t *testing.T) {
	market := v1.NewMockMarketPay()
	order, err := market.CreateOrder("name", "value")
	assert.NilError(t, err)
	assert.Assert(t, !market.IsPaid(*order))
	market.SetOrderStatus(v1.PUBLISHED)
	assert.Assert(t, market.IsPaid(*order))
}
