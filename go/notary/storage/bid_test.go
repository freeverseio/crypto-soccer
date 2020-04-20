package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestBidInsert(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auction := storage.NewAuction()
	auction.ID = "0"
	assert.NilError(t, auction.Insert(tx))

	bid := storage.NewBid()
	bid.AuctionID = auction.ID
	assert.NilError(t, bid.Insert(tx))
}

func TestBidsByAuctionID(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auction := storage.NewAuction()
	auction.ID = "0"
	assert.NilError(t, auction.Insert(tx))

	bid := storage.NewBid()
	bid.AuctionID = auction.ID
	assert.NilError(t, bid.Insert(tx))
	bid.ExtraPrice = 10
	assert.NilError(t, bid.Insert(tx))

	bids, err := storage.BidsByAuctionID(tx, auction.ID)
	assert.NilError(t, err)
	assert.Equal(t, len(bids), 2)

	auction.ID = "1"
	assert.NilError(t, auction.Insert(tx))

	bid = storage.NewBid()
	bid.AuctionID = auction.ID
	assert.NilError(t, bid.Insert(tx))

	bids, err = storage.BidsByAuctionID(tx, auction.ID)
	assert.NilError(t, err)
	assert.Equal(t, len(bids), 1)

}
