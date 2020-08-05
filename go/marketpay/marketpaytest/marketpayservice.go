package marketpaytest

import (
	"fmt"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/marketpay"

	"gotest.tools/assert"
)

func TestMarketPayService(t *testing.T, mp marketpay.MarketPayService) {
	t.Run("TestCreateOrder", func(t *testing.T) {
		name := "pippo"
		value := "134.10"
		order, err := mp.CreateOrder(name, value)
		assert.NilError(t, err)
		assert.Equal(t, order.Status, "DRAFT")
		assert.Equal(t, order.Amount, float64(134.1))
	})
	t.Run("TestCreateOrder2", func(t *testing.T) {
		auctionPrice := 4101
		extraPrice := 0
		price := fmt.Sprintf("%.2f", float64(auctionPrice+extraPrice)/100.0)
		name := "Freeverse Player transaction"
		order, err := mp.CreateOrder(name, price)
		assert.NilError(t, err)
		assert.Equal(t, order.Amount, float64(41.01))
		assert.Equal(t, order.Status, "DRAFT")
	})
	t.Run("TestGetOrder", func(t *testing.T) {
		name := "pippo"
		value := "134.10"
		order, err := mp.CreateOrder(name, value)
		assert.NilError(t, err)
		order1, err := mp.GetOrder(order.TrusteeShortlink.Hash)
		assert.NilError(t, err)
		assert.Equal(t, order.Name, order1.Name)
		_, err = mp.ValidateOrder(order.TrusteeShortlink.Hash)
		assert.NilError(t, err)
	})
	t.Run("TestIsPaid", func(t *testing.T) {
		name := "pippo"
		value := "134.10"
		order, err := mp.CreateOrder(name, value)
		if err != nil {
			t.Fatal(err)
		}
		isPaid := mp.IsPaid(*order)
		if isPaid {
			t.Fatal("Should not be paid")
		}
	})
}
