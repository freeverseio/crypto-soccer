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

func TestStorageHistoryUpdateChangedAuction(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.DB().Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auction := storage.NewAuction()
	assert.NilError(t, service.AuctionInsert(tx, *auction))
	auction.State = storage.AuctionCancelled
	assert.NilError(t, service.AuctionUpdate(tx, *auction))

	var count int
	tx.QueryRow("SELECT count(*) FROM auctions_histories;").Scan(&count)
	assert.Equal(t, count, 2)
}

func TestStorageHistoryUpdateUnexistentAuction(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.DB().Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auction := storage.NewAuction()
	assert.NilError(t, service.AuctionUpdate(tx, *auction))
}

func TestStorageHistoryUpdateUnchangedBid(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.DB().Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auction := storage.NewAuction()
	assert.NilError(t, service.AuctionInsert(tx, *auction))
	bid := storage.NewBid()
	bid.AuctionID = auction.ID
	assert.NilError(t, service.BidInsert(tx, *bid))
	assert.NilError(t, service.BidUpdate(tx, *bid))

	var count int
	tx.QueryRow("SELECT count(*) FROM bids_histories;").Scan(&count)
	assert.Equal(t, count, 1)
}

func TestStorageHistoryUpdateChangedBid(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.DB().Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auction := storage.NewAuction()
	assert.NilError(t, service.AuctionInsert(tx, *auction))
	bid := storage.NewBid()
	bid.AuctionID = auction.ID
	assert.NilError(t, service.BidInsert(tx, *bid))
	bid.State = storage.BidPaid
	assert.NilError(t, service.BidUpdate(tx, *bid))

	var count int
	tx.QueryRow("SELECT count(*) FROM bids_histories;").Scan(&count)
	assert.Equal(t, count, 2)
}

func TestStorageHistoryUpdateUnChangedPlaystoreOrder(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.DB().Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	order := storage.NewPlaystoreOrder()
	assert.NilError(t, service.PlayStoreInsert(tx, *order))
	assert.NilError(t, service.PlayStoreUpdateState(tx, *order))

	var count int
	tx.QueryRow("SELECT count(*) FROM playstore_orders_histories;").Scan(&count)
	assert.Equal(t, count, 1)
}

func TestStorageHistoryUpdateChangedPlaystoreOrder(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.DB().Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	order := storage.NewPlaystoreOrder()
	assert.NilError(t, service.PlayStoreInsert(tx, *order))
	order.PlayerId = "234234234"
	assert.NilError(t, service.PlayStoreUpdateState(tx, *order))

	var count int
	tx.QueryRow("SELECT count(*) FROM playstore_orders_histories;").Scan(&count)
	assert.Equal(t, count, 2)
}

func TestStorageHistoryUpdateUnChangedOffer(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.DB().Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	order := storage.NewOffer()
	assert.NilError(t, service.OfferInsert(tx, *order))
	assert.NilError(t, service.OfferUpdate(tx, *order))

	var count int
	tx.QueryRow("SELECT count(*) FROM offers_histories;").Scan(&count)
	assert.Equal(t, count, 1)
}

func TestStorageHistoryUpdateChangedOffer(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.DB().Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auction := storage.NewAuction()
	assert.NilError(t, service.AuctionInsert(tx, *auction))

	order := storage.NewOffer()
	assert.NilError(t, service.OfferInsert(tx, *order))
	order.AuctionID = auction.ID
	order.Seller = "me"
	assert.NilError(t, service.OfferUpdate(tx, *order))

	var count int
	tx.QueryRow("SELECT count(*) FROM offers_histories;").Scan(&count)
	assert.Equal(t, count, 2)
}
