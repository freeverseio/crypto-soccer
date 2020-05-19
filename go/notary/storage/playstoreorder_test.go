package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestPlaystoreOrderCreate(t *testing.T) {
	order := storage.NewPlaystoreOrder()
	assert.Equal(t, order.State, storage.PlaystoreOrderPending)
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
	order.State = storage.PlaystoreOrderFailed
	order.StateExtra = "prova"
	assert.NilError(t, order.Insert(tx))

	result, err := storage.PlaystoreOrderByOrderId(tx, order.OrderId)
	assert.NilError(t, err)
	assert.Assert(t, result != nil)
	assert.Equal(t, *result, *order)
}
