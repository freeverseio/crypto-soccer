package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestPlaystoreOrderCreate(t *testing.T) {
	orderId := "ciao"
	order := storage.NewPlaystoreOrder(orderId)
	assert.Equal(t, order.OrderId, orderId)
	assert.Equal(t, order.State, storage.PlaystoreOrderPending)
	assert.Equal(t, order.StateExtra, "")
}

func TestPlaystoreOrderInsert(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	orderId := "ciao"
	order := storage.NewPlaystoreOrder(orderId)
	order.State = storage.PlaystoreOrderFailed
	order.StateExtra = "prova"
	assert.NilError(t, order.Insert(tx))

	result, err := storage.PlaystoreOrderByOrderId(tx, orderId)
	assert.NilError(t, err)
	assert.Equal(t, *result, *order)
}
