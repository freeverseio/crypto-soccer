package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestAuctionByIDUnexistent(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auction, err := storage.AuctionByID(tx, "4343")
	assert.NilError(t, err)
	assert.Assert(t, auction == nil)
}

func TestAuctionInsert(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auction := storage.NewAuction()
	auction.ID = "ciao"
	auction.Rnd = 4
	auction.PlayerID = "3"
	auction.CurrencyID = 3
	auction.Price = 3
	auction.ValidUntil = 3
	auction.Signature = "3"
	auction.State = storage.AuctionStarted
	auction.StateExtra = "3"
	auction.PaymentURL = "3"
	auction.Seller = "3"
	assert.NilError(t, auction.Insert(tx))

	result, err := storage.AuctionByID(tx, auction.ID)
	assert.NilError(t, err)
	assert.Equal(t, *result, *auction)
}

func TestPendingAuctions(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auction := storage.NewAuction()
	auction.ID = "ciao"
	auction.State = storage.AuctionStarted
	assert.NilError(t, auction.Insert(tx))
	result, err := storage.PendingAuctions(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(result), 1)

	auction.ID = "ciao1"
	auction.State = storage.AuctionAssetFrozen
	assert.NilError(t, auction.Insert(tx))
	result, err = storage.PendingAuctions(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(result), 2)

	auction.ID = "ciao2"
	auction.State = storage.AuctionPaying
	assert.NilError(t, auction.Insert(tx))
	result, err = storage.PendingAuctions(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(result), 3)

	auction.ID = "ciao3"
	auction.State = storage.AuctionWithdrableBySeller
	assert.NilError(t, auction.Insert(tx))
	result, err = storage.PendingAuctions(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(result), 4)

	auction.ID = "ciao4"
	auction.State = storage.AuctionWithdrableByBuyer
	assert.NilError(t, auction.Insert(tx))
	result, err = storage.PendingAuctions(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(result), 5)

	auction.ID = "ciao5"
	auction.State = storage.AuctionFailed
	assert.NilError(t, auction.Insert(tx))
	result, err = storage.PendingAuctions(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(result), 5)

	auction.ID = "ciao6"
	auction.State = storage.AuctionEnded
	assert.NilError(t, auction.Insert(tx))
	result, err = storage.PendingAuctions(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(result), 5)

	auction.ID = "ciao7"
	auction.State = storage.AuctionCancelled
	assert.NilError(t, auction.Insert(tx))
	result, err = storage.PendingAuctions(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(result), 5)

}

func TestAuctionUpdate(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auction := storage.NewAuction()
	auction.ID = "ciao"
	auction.State = storage.AuctionStarted
	auction.StateExtra = "priva"
	assert.NilError(t, auction.Insert(tx))
	result, err := storage.AuctionByID(tx, auction.ID)
	assert.NilError(t, err)
	assert.Equal(t, result.State, storage.AuctionStarted)
	assert.Equal(t, result.StateExtra, "priva")

	auction.State = storage.AuctionCancelled
	auction.StateExtra = "privato"
	auction.PaymentURL = "http"
	assert.NilError(t, auction.Update(tx))

	result, err = storage.AuctionByID(tx, auction.ID)
	assert.NilError(t, err)
	assert.Equal(t, *result, *auction)
}
