package storagetest

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func testPlaystoreOrderServiceInterface(t *testing.T, service storage.StorageService) {
	t.Run("insert", func(t *testing.T) {
		assert.NilError(t, service.Begin())
		defer service.Rollback()

		order := storage.NewPlaystoreOrder()
		order.PurchaseToken = "ciao"
		order.PackageName = "dsd"
		order.ProductId = "444"
		order.OrderId = "fdrd"
		order.PlayerId = "4"
		order.TeamId = "pippo"
		order.State = storage.PlaystoreOrderFailed
		order.StateExtra = "prova"
		order.Signature = "erere"

		assert.NilError(t, service.PlayStoreInsert(*order))
		result, err := service.PlayStoreOrder(order.OrderId)
		assert.NilError(t, err)
		assert.Assert(t, result != nil)
		assert.Equal(t, *result, *order)
	})

	t.Run("pending orders", func(t *testing.T) {
		assert.NilError(t, service.Begin())
		defer service.Rollback()

		order := storage.NewPlaystoreOrder()
		order.PurchaseToken = "ciao1"
		order.PackageName = "dsd"
		order.ProductId = "444"
		order.OrderId = "fdrd"
		order.PlayerId = "4"
		order.TeamId = "pippo"
		order.State = storage.PlaystoreOrderFailed
		order.StateExtra = "prova"
		order.Signature = "erere"
		assert.NilError(t, service.PlayStoreInsert(*order))

		orders, err := service.PlayStorePendingOrders()
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 0)

		order.PurchaseToken = "43d"
		order.State = storage.PlaystoreOrderOpen
		assert.NilError(t, service.PlayStoreInsert(*order))

		orders, err = service.PlayStorePendingOrders()
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 1)

		order.PurchaseToken = "43d1"
		order.State = storage.PlaystoreOrderAcknowledged
		assert.NilError(t, service.PlayStoreInsert(*order))

		orders, err = service.PlayStorePendingOrders()
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 2)

		order.PurchaseToken = "43d2"
		order.State = storage.PlaystoreOrderComplete
		assert.NilError(t, service.PlayStoreInsert(*order))

		orders, err = service.PlayStorePendingOrders()
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 2)
	})

	t.Run("update state", func(t *testing.T) {
		assert.NilError(t, service.Begin())
		defer service.Rollback()

		order := storage.NewPlaystoreOrder()
		order.PurchaseToken = "ciao"
		order.PackageName = "dsd"
		order.ProductId = "444"
		order.OrderId = "fdrd"
		order.PlayerId = "4"
		order.TeamId = "pippo"
		order.State = storage.PlaystoreOrderFailed
		order.StateExtra = "prova"
		order.Signature = "erere"

		assert.NilError(t, service.PlayStoreInsert(*order))
		assert.NilError(t, service.PlayStoreUpdateState(*order))

		order.State = storage.PlaystoreOrderOpen
		order.StateExtra = "recdia"
		assert.NilError(t, service.PlayStoreUpdateState(*order))

		result, err := service.PlayStoreOrder(order.OrderId)
		assert.NilError(t, err)
		assert.Assert(t, result != nil)
		assert.Equal(t, result.State, order.State)
		assert.Equal(t, result.StateExtra, order.StateExtra)

	})

	t.Run("pending order by playerId", func(t *testing.T) {
		assert.NilError(t, service.Begin())
		defer service.Rollback()

		order := storage.NewPlaystoreOrder()
		order.PurchaseToken = "ciao12"
		order.PackageName = "dsd"
		order.ProductId = "444"
		order.OrderId = "fdrd"
		order.PlayerId = "4343534"
		order.TeamId = "pippo"
		order.State = storage.PlaystoreOrderFailed
		order.StateExtra = "prova"
		order.Signature = "erere"

		orders, err := service.PlayStorePendingOrdersByPlayerId(order.PlayerId)
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 0)

		assert.NilError(t, service.PlayStoreInsert(*order))
		orders, err = service.PlayStorePendingOrdersByPlayerId(order.PlayerId)
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 0)

		order.PurchaseToken = "ciao432423"
		order.State = storage.PlaystoreOrderComplete
		assert.NilError(t, service.PlayStoreInsert(*order))
		orders, err = service.PlayStorePendingOrdersByPlayerId(order.PlayerId)
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 0)

		order.PurchaseToken = "ciao4324233"
		order.State = storage.PlaystoreOrderAcknowledged
		assert.NilError(t, service.PlayStoreInsert(*order))
		orders, err = service.PlayStorePendingOrdersByPlayerId(order.PlayerId)
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 1)
	})
}
