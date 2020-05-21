package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestPlaystoreOrderCreate(t *testing.T) {
	order := storage.NewPlaystoreOrder()
	assert.Equal(t, order.State, storage.PlaystoreOrderOpen)
}

func TestPlaystoreOrderInsert(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	order := storage.NewPlaystoreOrder()
	order.OrderId = "ciao"
	order.PackageName = "dsd"
	order.ProductId = "444"
	order.PurchaseToken = "fdrd"
	order.PlayerId = "4"
	order.TeamId = "pippo"
	order.State = storage.PlaystoreOrderFailed
	order.StateExtra = "prova"
	order.Signature = "erere"
	assert.NilError(t, order.Insert(tx))

	result, err := storage.PlaystoreOrderByOrderId(tx, order.OrderId)
	assert.NilError(t, err)
	assert.Assert(t, result != nil)
	assert.Equal(t, *result, *order)
}

func TestPlaystoreOrderOpenOrders(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	order := storage.NewPlaystoreOrder()
	order.OrderId = "ciao"
	order.PackageName = "dsd"
	order.ProductId = "444"
	order.PurchaseToken = "fdrd"
	order.PlayerId = "4"
	order.TeamId = "pippo"
	order.State = storage.PlaystoreOrderFailed
	order.StateExtra = "prova"
	order.Signature = "erere"
	assert.NilError(t, order.Insert(tx))

	orders, err := storage.PendingPlaystoreOrders(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(orders), 0)

	order.OrderId = "43d"
	order.State = storage.PlaystoreOrderOpen
	assert.NilError(t, order.Insert(tx))

	orders, err = storage.PendingPlaystoreOrders(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(orders), 1)

	order.OrderId = "43d1"
	order.State = storage.PlaystoreOrderAcknowledged
	assert.NilError(t, order.Insert(tx))

	orders, err = storage.PendingPlaystoreOrders(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(orders), 2)

	order.OrderId = "43d2"
	order.State = storage.PlaystoreOrderComplete
	assert.NilError(t, order.Insert(tx))

	orders, err = storage.PendingPlaystoreOrders(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(orders), 2)
}

func TestPlaystoreOrderUpdateState(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	order := storage.NewPlaystoreOrder()
	order.OrderId = "ciao"
	order.PackageName = "dsd"
	order.ProductId = "444"
	order.PurchaseToken = "fdrd"
	order.PlayerId = "4"
	order.TeamId = "pippo"
	order.State = storage.PlaystoreOrderFailed
	order.StateExtra = "prova"
	order.Signature = "erere"
	assert.NilError(t, order.Insert(tx))

	order.State = storage.PlaystoreOrderOpen
	order.StateExtra = "recdia"
	assert.NilError(t, order.UpdateState(tx))

	result, err := storage.PlaystoreOrderByOrderId(tx, order.OrderId)
	assert.NilError(t, err)
	assert.Equal(t, result.State, order.State)
	assert.Equal(t, result.StateExtra, order.StateExtra)
}
