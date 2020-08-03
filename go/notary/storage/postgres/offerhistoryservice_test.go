package postgres_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	"gotest.tools/assert"
)

func TestOfferHistoryServiceInterface(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	service := postgres.NewOfferHistoryService(tx)
	auctionService := postgres.NewAuctionService(tx)
	storage.TestOfferServiceInterface(t, service, auctionService)
}
