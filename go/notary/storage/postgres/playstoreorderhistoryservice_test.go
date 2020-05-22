package postgres_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	"gotest.tools/assert"
)

func TestPlaystoreHistoryOrder(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	service := postgres.NewPlaystoreOrderService(tx)
	historyService := postgres.NewPlaystoreOrderHistoryService(tx, service)
	storage.TestPlaystoreOrderServiceInterface(t, historyService)
}
