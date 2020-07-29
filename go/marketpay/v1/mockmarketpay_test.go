package v1_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/marketpay"
	"gotest.tools/assert"
)

func TestMockMarketPaySetStatus(t *testing.T) {
	market := marketpay.NewMockMarketPay()
	order, err := market.CreateOrder("name", "value")
	assert.NilError(t, err)
	assert.Assert(t, !market.IsPaid(*order))
	market.SetOrderStatus(marketpay.PUBLISHED)
	assert.Assert(t, market.IsPaid(*order))
}
