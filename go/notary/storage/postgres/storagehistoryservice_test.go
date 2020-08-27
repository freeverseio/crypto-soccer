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

	assert.NilError(t, service.Begin())
	defer service.Rollback()

	assert.Equal(t, service.AuctionsHistoriesCount(), 0)

	auction := storage.NewAuction()
	assert.NilError(t, service.AuctionInsert(*auction))

	assert.Equal(t, service.AuctionsHistoriesCount(), 1)
}

func TestStorageHistoryUpdateUnchangedAuction(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	assert.NilError(t, service.Begin())
	defer service.Rollback()

	auction := storage.NewAuction()
	assert.NilError(t, service.AuctionInsert(*auction))
	assert.NilError(t, service.AuctionUpdate(*auction))

	assert.Equal(t, service.AuctionsHistoriesCount(), 1)
}

func TestStorageHistoryUpdateChangedAuction(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	assert.NilError(t, service.Begin())
	defer service.Rollback()

	auction := storage.NewAuction()
	assert.NilError(t, service.AuctionInsert(*auction))
	auction.State = storage.AuctionCancelled
	assert.NilError(t, service.AuctionUpdate(*auction))

	assert.Equal(t, service.AuctionsHistoriesCount(), 2)
}

func TestStorageHistoryUpdateUnexistentAuction(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	assert.NilError(t, service.Begin())
	defer service.Rollback()

	auction := storage.NewAuction()
	assert.NilError(t, service.AuctionUpdate(*auction))
}

func TestStorageHistoryUpdateUnchangedBid(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	assert.NilError(t, service.Begin())
	defer service.Rollback()

	auction := storage.NewAuction()
	assert.NilError(t, service.AuctionInsert(*auction))
	bid := storage.NewBid()
	bid.AuctionID = auction.ID
	assert.NilError(t, service.BidInsert(*bid))
	assert.NilError(t, service.BidUpdate(*bid))

	assert.Equal(t, service.BidsHistoriesCount(), 1)
}

func TestStorageHistoryUpdateChangedBid(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	assert.NilError(t, service.Begin())
	defer service.Rollback()

	auction := storage.NewAuction()
	assert.NilError(t, service.AuctionInsert(*auction))
	bid := storage.NewBid()
	bid.AuctionID = auction.ID
	assert.NilError(t, service.BidInsert(*bid))
	bid.State = storage.BidPaid
	assert.NilError(t, service.BidUpdate(*bid))

	assert.Equal(t, service.BidsHistoriesCount(), 2)
}

func TestStorageHistoryUpdateUnChangedPlaystoreOrder(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	assert.NilError(t, service.Begin())
	defer service.Rollback()

	order := storage.NewPlaystoreOrder()
	assert.NilError(t, service.PlayStoreInsert(*order))
	assert.NilError(t, service.PlayStoreUpdateState(*order))

	assert.Equal(t, service.PlaystoreHistoriesCount(), 1)
}

func TestStorageHistoryUpdateChangedPlaystoreOrder(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	assert.NilError(t, service.Begin())
	defer service.Rollback()

	order := storage.NewPlaystoreOrder()
	assert.NilError(t, service.PlayStoreInsert(*order))
	order.PlayerId = "234234234"
	assert.NilError(t, service.PlayStoreUpdateState(*order))

	assert.Equal(t, service.PlaystoreHistoriesCount(), 2)
}

func TestStorageHistoryUpdateUnChangedOffer(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	assert.NilError(t, service.Begin())
	defer service.Rollback()

	order := storage.NewOffer()
	assert.NilError(t, service.OfferInsert(*order))
	assert.NilError(t, service.OfferUpdate(*order))

	assert.Equal(t, service.OffersHistoriesCount(), 1)
}

func TestStorageHistoryUpdateChangedOffer(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	assert.NilError(t, service.Begin())
	defer service.Rollback()

	auction := storage.NewAuction()
	assert.NilError(t, service.AuctionInsert(*auction))

	order := storage.NewOffer()
	assert.NilError(t, service.OfferInsert(*order))
	order.AuctionID = auction.ID
	order.Seller = "me"
	assert.NilError(t, service.OfferUpdate(*order))

	assert.Equal(t, service.OffersHistoriesCount(), 2)
}
