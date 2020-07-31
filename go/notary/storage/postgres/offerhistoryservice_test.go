package postgres_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/storagetest"
	"gotest.tools/assert"
)

func TestOfferHistoryServiceInterface(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	service := postgres.NewOfferHistoryService(tx)
	auctionService := postgres.NewAuctionService(tx)
	storagetest.TestOfferServiceInterface(t, service, auctionService)
}
