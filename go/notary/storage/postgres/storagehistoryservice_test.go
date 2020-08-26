package postgres_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/storagetest"
	"gotest.tools/assert"
)

func TestStorageHistoryServiceStart(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)
	storagetest.TestStorageService(t, service)
}

func TestStorageHistoryInsertAuction(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.DB().Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	var count int
	tx.QueryRow("SELECT count(*) FROM auctions_histories;").Scan(&count)
	assert.Equal(t, count, 0)

	auction := storage.NewAuction()
	assert.NilError(t, service.AuctionInsert(tx, *auction))

	tx.QueryRow("SELECT count(*) FROM auctions_histories;").Scan(&count)
	assert.Equal(t, count, 1)
}

func TestStorageHistoryUpdateUnchangedAuction(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.DB().Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auction := storage.NewAuction()
	assert.NilError(t, service.AuctionInsert(tx, *auction))
	assert.NilError(t, service.AuctionUpdate(tx, *auction))

	var count int
	tx.QueryRow("SELECT count(*) FROM auctions_histories;").Scan(&count)
	assert.Equal(t, count, 1)
}

func TestStorageHistoryUpdateUnexistentAuction(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.DB().Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auction := storage.NewAuction()
	assert.NilError(t, service.AuctionUpdate(tx, *auction))
}
