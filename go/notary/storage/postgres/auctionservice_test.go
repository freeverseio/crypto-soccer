package postgres_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	"gotest.tools/assert"
)

func TestAuctionServiceInterface(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	service := postgres.NewAuctionService(tx)
	storage.TestAuctionServiceInterface(t, service)
}
