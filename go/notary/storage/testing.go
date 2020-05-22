package storage

import (
	"testing"

	"gotest.tools/assert"
)

func TestPlaystoreOrderServiceInterface(t *testing.T, service PlaystoreOrderService) {
	order := NewPlaystoreOrder()
	order.OrderId = "ciao"
	order.PackageName = "dsd"
	order.ProductId = "444"
	order.PurchaseToken = "fdrd"
	order.PlayerId = "4"
	order.TeamId = "pippo"
	order.State = PlaystoreOrderFailed
	order.StateExtra = "prova"
	order.Signature = "erere"

	assert.NilError(t, service.Insert(*order))
	result, err := service.Order(order.OrderId)
	assert.NilError(t, err)
	assert.Assert(t, result != nil)
	assert.Equal(t, *result, *order)
}
