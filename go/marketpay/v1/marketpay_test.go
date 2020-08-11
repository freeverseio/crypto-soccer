package v1_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/marketpay/marketpaytest"

	v1 "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
)

func TestMarketPayService(t *testing.T) {
	service := v1.New("")
	marketpaytest.TestMarketPayService(t, service)
}
