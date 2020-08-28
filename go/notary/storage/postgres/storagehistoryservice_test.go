package postgres_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/storagetest"
	"gotest.tools/assert"
)

func TestStorageHistoryServiceStart(t *testing.T) {
	tx := postgres.NewStorageHistoryService(db)
	storagetest.TestStorageService(t, tx)
}

func TestStorageHistoryInsertAuction(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	switch tx := tx.(type) {
	default:
		t.Errorf("wrong type %T", tx)
	case *postgres.StorageHistoryTx:
		assert.Equal(t, tx.AuctionsHistoriesCount(), 0)

		auction := storage.NewAuction()
		assert.NilError(t, tx.AuctionInsert(*auction))

		assert.Equal(t, tx.AuctionsHistoriesCount(), 1)
	}
}

func TestStorageHistoryUpdateUnchangedAuction(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	switch tx := tx.(type) {
	default:
		t.Errorf("wrong type %T", tx)
	case *postgres.StorageHistoryTx:

		auction := storage.NewAuction()
		assert.NilError(t, tx.AuctionInsert(*auction))
		assert.NilError(t, tx.AuctionUpdate(*auction))

		assert.Equal(t, tx.AuctionsHistoriesCount(), 1)
	}
}

func TestStorageHistoryUpdateChangedAuction(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	switch tx := tx.(type) {
	default:
		t.Errorf("wrong type %T", tx)
	case *postgres.StorageHistoryTx:

		auction := storage.NewAuction()
		assert.NilError(t, tx.AuctionInsert(*auction))
		auction.State = storage.AuctionCancelled
		assert.NilError(t, tx.AuctionUpdate(*auction))

		assert.Equal(t, tx.AuctionsHistoriesCount(), 2)
	}
}

func TestStorageHistoryUpdateUnexistentAuction(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auction := storage.NewAuction()
	assert.NilError(t, tx.AuctionUpdate(*auction))
}

func TestStorageHistoryUpdateUnchangedBid(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	switch tx := tx.(type) {
	default:
		t.Errorf("wrong type %T", tx)
	case *postgres.StorageHistoryTx:
		auction := storage.NewAuction()
		assert.NilError(t, tx.AuctionInsert(*auction))
		bid := storage.NewBid()
		bid.AuctionID = auction.ID
		assert.NilError(t, tx.BidInsert(*bid))
		assert.NilError(t, tx.BidUpdate(*bid))

		assert.Equal(t, tx.BidsHistoriesCount(), 1)
	}
}

func TestStorageHistoryUpdateChangedBid(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	switch tx := tx.(type) {
	default:
		t.Errorf("wrong type %T", tx)
	case *postgres.StorageHistoryTx:
		auction := storage.NewAuction()
		assert.NilError(t, tx.AuctionInsert(*auction))
		bid := storage.NewBid()
		bid.AuctionID = auction.ID
		assert.NilError(t, tx.BidInsert(*bid))
		bid.State = storage.BidPaid
		assert.NilError(t, tx.BidUpdate(*bid))
		assert.Equal(t, tx.BidsHistoriesCount(), 2)
	}
}

func TestStorageHistoryUpdateUnChangedPlaystoreOrder(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	switch tx := tx.(type) {
	default:
		t.Errorf("wrong type %T", tx)
	case *postgres.StorageHistoryTx:
		order := storage.NewPlaystoreOrder()
		assert.NilError(t, tx.PlayStoreInsert(*order))
		assert.NilError(t, tx.PlayStoreUpdateState(*order))

		assert.Equal(t, tx.PlaystoreHistoriesCount(), 1)
	}
}

func TestStorageHistoryUpdateChangedPlaystoreOrder(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	switch tx := tx.(type) {
	default:
		t.Errorf("wrong type %T", tx)
	case *postgres.StorageHistoryTx:
		order := storage.NewPlaystoreOrder()
		assert.NilError(t, tx.PlayStoreInsert(*order))
		order.PlayerId = "234234234"
		assert.NilError(t, tx.PlayStoreUpdateState(*order))

		assert.Equal(t, tx.PlaystoreHistoriesCount(), 2)
	}
}

func TestStorageHistoryUpdateUnChangedOffer(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	switch tx := tx.(type) {
	default:
		t.Errorf("wrong type %T", tx)
	case *postgres.StorageHistoryTx:
		order := storage.NewOffer()
		assert.NilError(t, tx.OfferInsert(*order))
		assert.NilError(t, tx.OfferUpdate(*order))
		assert.Equal(t, tx.OffersHistoriesCount(), 1)
	}
}

func TestStorageHistoryUpdateChangedOffer(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)

	tx, err := service.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	switch tx := tx.(type) {
	default:
		t.Errorf("wrong type %T", tx)
	case *postgres.StorageHistoryTx:
		auction := storage.NewAuction()
		assert.NilError(t, tx.AuctionInsert(*auction))

		order := storage.NewOffer()
		assert.NilError(t, tx.OfferInsert(*order))
		order.AuctionID = auction.ID
		order.Seller = "me"
		assert.NilError(t, tx.OfferUpdate(*order))

		assert.Equal(t, tx.OffersHistoriesCount(), 2)
	}
}
