package storage

import (
	"testing"

	"gotest.tools/assert"
)

func TestPlaystoreOrderServiceInterface(t *testing.T, service PlaystoreOrderService) {
	t.Run("insert", func(t *testing.T) {
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
	})

	t.Run("pending orders", func(t *testing.T) {
		order := NewPlaystoreOrder()
		order.OrderId = "ciao1"
		order.PackageName = "dsd"
		order.ProductId = "444"
		order.PurchaseToken = "fdrd"
		order.PlayerId = "4"
		order.TeamId = "pippo"
		order.State = PlaystoreOrderFailed
		order.StateExtra = "prova"
		order.Signature = "erere"
		assert.NilError(t, service.Insert(*order))

		orders, err := service.PendingOrders()
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 0)

		order.OrderId = "43d"
		order.State = PlaystoreOrderOpen
		assert.NilError(t, service.Insert(*order))

		orders, err = service.PendingOrders()
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 1)

		order.OrderId = "43d1"
		order.State = PlaystoreOrderAcknowledged
		assert.NilError(t, service.Insert(*order))

		orders, err = service.PendingOrders()
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 2)

		order.OrderId = "43d2"
		order.State = PlaystoreOrderComplete
		assert.NilError(t, service.Insert(*order))

		orders, err = service.PendingOrders()
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 2)
	})

	t.Run("update state", func(t *testing.T) {
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
		assert.NilError(t, service.UpdateState(*order))

		order.State = PlaystoreOrderOpen
		order.StateExtra = "recdia"
		assert.NilError(t, service.UpdateState(*order))

		result, err := service.Order(order.OrderId)
		assert.NilError(t, err)
		assert.Equal(t, result.State, order.State)
		assert.Equal(t, result.StateExtra, order.StateExtra)

	})
}
