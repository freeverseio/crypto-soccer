package storagetest

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func testPlaystoreOrderServiceInterface(t *testing.T, service storage.StorageService) {
	t.Run("insert", func(t *testing.T) {
		tx, err := service.DB().Begin()
		assert.NilError(t, err)
		defer tx.Rollback()

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

		assert.NilError(t, service.PlayStoreInsert(tx, *order))
		result, err := service.PlayStoreOrder(tx, order.OrderId)
		assert.NilError(t, err)
		assert.Assert(t, result != nil)
		assert.Equal(t, *result, *order)
	})

	t.Run("pending orders", func(t *testing.T) {
		tx, err := service.DB().Begin()
		assert.NilError(t, err)
		defer tx.Rollback()

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
		assert.NilError(t, service.PlayStoreInsert(tx, *order))

		orders, err := service.PlayStorePendingOrders(tx)
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 0)

		order.PurchaseToken = "43d"
		order.State = storage.PlaystoreOrderOpen
		assert.NilError(t, service.PlayStoreInsert(tx, *order))

		orders, err = service.PlayStorePendingOrders(tx)
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 1)

		order.PurchaseToken = "43d1"
		order.State = storage.PlaystoreOrderAcknowledged
		assert.NilError(t, service.PlayStoreInsert(tx, *order))

		orders, err = service.PlayStorePendingOrders(tx)
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 2)

		order.PurchaseToken = "43d2"
		order.State = storage.PlaystoreOrderComplete
		assert.NilError(t, service.PlayStoreInsert(tx, *order))

		orders, err = service.PlayStorePendingOrders(tx)
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 2)
	})

	t.Run("update state", func(t *testing.T) {
		tx, err := service.DB().Begin()
		assert.NilError(t, err)
		defer tx.Rollback()

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

		assert.NilError(t, service.PlayStoreInsert(tx, *order))
		assert.NilError(t, service.PlayStoreUpdateState(tx, *order))

		order.State = storage.PlaystoreOrderOpen
		order.StateExtra = "recdia"
		assert.NilError(t, service.PlayStoreUpdateState(tx, *order))

		result, err := service.PlayStoreOrder(tx, order.OrderId)
		assert.NilError(t, err)
		assert.Assert(t, result != nil)
		assert.Equal(t, result.State, order.State)
		assert.Equal(t, result.StateExtra, order.StateExtra)

	})

	t.Run("pending order by playerId", func(t *testing.T) {
		tx, err := service.DB().Begin()
		assert.NilError(t, err)
		defer tx.Rollback()

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

		orders, err := service.PlayStorePendingOrdersByPlayerId(tx, order.PlayerId)
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 0)

		assert.NilError(t, service.PlayStoreInsert(tx, *order))
		orders, err = service.PlayStorePendingOrdersByPlayerId(tx, order.PlayerId)
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 0)

		order.PurchaseToken = "ciao432423"
		order.State = storage.PlaystoreOrderComplete
		assert.NilError(t, service.PlayStoreInsert(tx, *order))
		orders, err = service.PlayStorePendingOrdersByPlayerId(tx, order.PlayerId)
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 0)

		order.PurchaseToken = "ciao4324233"
		order.State = storage.PlaystoreOrderAcknowledged
		assert.NilError(t, service.PlayStoreInsert(tx, *order))
		orders, err = service.PlayStorePendingOrdersByPlayerId(tx, order.PlayerId)
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 1)
	})
}
