package v1_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/marketpay/marketpaytest"
	"gotest.tools/assert"

	v1 "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
)

func TestMarketPayService(t *testing.T) {
	service := v1.New("")
	marketpaytest.TestMarketPayService(t, service)
}

func TestMarketPayServiceCreateOrderWithSameParams(t *testing.T) {
	service := v1.NewSandbox()
	name := "ciao"
	value := "10"
	// any time an order is created it returns different hashes even if we use same params
	order0, err := service.CreateOrder(name, value)
	assert.NilError(t, err)
	order1, err := service.CreateOrder(name, value)
	assert.NilError(t, err)
	order2, err := service.CreateOrder(name+"1", value+"1")
	assert.NilError(t, err)
	assert.Equal(t, order1.TrusteeShortlink.Hash == order0.TrusteeShortlink.Hash, false)
	assert.Equal(t, order2.TrusteeShortlink.Hash == order0.TrusteeShortlink.Hash, false)
}
