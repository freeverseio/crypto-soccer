package storagetest

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func testAuctionServiceInterface(t *testing.T, service storage.StorageService) {
	t.Run("TestAuctionByIDUnexistent", func(t *testing.T) {
		tx, err := service.DB().Begin()
		assert.NilError(t, err)
		defer tx.Rollback()

		auction, err := service.Auction(tx, "4343")
		assert.NilError(t, err)
		assert.Assert(t, auction == nil)
	})

	t.Run("TestAuctionInsert", func(t *testing.T) {
		tx, err := service.DB().Begin()
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
		assert.NilError(t, service.AuctionInsert(tx, *auction))

		result, err := service.Auction(tx, auction.ID)
		assert.NilError(t, err)
		assert.Equal(t, *result, *auction)
	})

	t.Run("TestAuctionPendingAuctions", func(t *testing.T) {
		tx, err := service.DB().Begin()
		assert.NilError(t, err)
		defer tx.Rollback()

		auction := storage.NewAuction()
		auction.ID = "ciao0"
		auction.State = storage.AuctionStarted
		assert.NilError(t, service.AuctionInsert(tx, *auction))
		result, err := service.AuctionPendingAuctions(tx)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 1)

		auction.ID = "ciao1"
		auction.State = storage.AuctionAssetFrozen
		assert.NilError(t, service.AuctionInsert(tx, *auction))
		result, err = service.AuctionPendingAuctions(tx)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 2)

		auction.ID = "ciao2"
		auction.State = storage.AuctionPaying
		assert.NilError(t, service.AuctionInsert(tx, *auction))
		result, err = service.AuctionPendingAuctions(tx)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 3)

		auction.ID = "ciao3"
		auction.State = storage.AuctionWithdrableBySeller
		assert.NilError(t, service.AuctionInsert(tx, *auction))
		result, err = service.AuctionPendingAuctions(tx)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 4)

		auction.ID = "ciao4"
		auction.State = storage.AuctionWithdrableByBuyer
		assert.NilError(t, service.AuctionInsert(tx, *auction))
		result, err = service.AuctionPendingAuctions(tx)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 5)

		auction.ID = "ciao5"
		auction.State = storage.AuctionFailed
		assert.NilError(t, service.AuctionInsert(tx, *auction))
		result, err = service.AuctionPendingAuctions(tx)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 5)

		auction.ID = "ciao6"
		auction.State = storage.AuctionEnded
		assert.NilError(t, service.AuctionInsert(tx, *auction))
		result, err = service.AuctionPendingAuctions(tx)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 5)

		auction.ID = "ciao7"
		auction.State = storage.AuctionCancelled
		assert.NilError(t, service.AuctionInsert(tx, *auction))
		result, err = service.AuctionPendingAuctions(tx)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 5)
	})

	t.Run("TestAuctionUpdate", func(t *testing.T) {
		tx, err := service.DB().Begin()
		assert.NilError(t, err)
		defer tx.Rollback()

		auction := storage.NewAuction()
		auction.ID = "ciao20"
		auction.State = storage.AuctionStarted
		auction.StateExtra = "priva"
		assert.NilError(t, service.AuctionInsert(tx, *auction))
		result, err := service.Auction(tx, auction.ID)
		assert.NilError(t, err)
		assert.Equal(t, result.State, storage.AuctionStarted)
		assert.Equal(t, result.StateExtra, "priva")

		auction.State = storage.AuctionCancelled
		auction.StateExtra = "privato"
		auction.PaymentURL = "http"
		assert.NilError(t, service.AuctionUpdate(tx, *auction))

		result, err = service.Auction(tx, auction.ID)
		assert.NilError(t, err)
		assert.Equal(t, *result, *auction)
	})

	t.Run("TestBid().Insert", func(t *testing.T) {
		tx, err := service.DB().Begin()
		assert.NilError(t, err)
		defer tx.Rollback()

		auction := storage.NewAuction()
		auction.ID = "0"
		assert.NilError(t, service.AuctionInsert(tx, *auction))

		bid := storage.NewBid()
		bid.AuctionID = auction.ID
		assert.NilError(t, service.BidInsert(tx, *bid))
	})

	t.Run("TestBidsByAuctionID", func(t *testing.T) {
		tx, err := service.DB().Begin()
		assert.NilError(t, err)
		defer tx.Rollback()

		auction := storage.NewAuction()
		auction.ID = "03"
		assert.NilError(t, service.AuctionInsert(tx, *auction))

		bid := storage.NewBid()
		bid.AuctionID = auction.ID
		assert.NilError(t, service.BidInsert(tx, *bid))
		bid.ExtraPrice = 10
		assert.NilError(t, service.BidInsert(tx, *bid))

		bids, err := service.Bids(tx, auction.ID)
		assert.NilError(t, err)
		assert.Equal(t, len(bids), 2)

		auction.ID = "1"
		assert.NilError(t, service.AuctionInsert(tx, *auction))

		bid = storage.NewBid()
		bid.AuctionID = auction.ID
		assert.NilError(t, service.BidInsert(tx, *bid))

		bids, err = service.Bids(tx, auction.ID)
		assert.NilError(t, err)
		assert.Equal(t, len(bids), 1)
	})

	t.Run("TestBidUpdate", func(t *testing.T) {
		tx, err := service.DB().Begin()
		assert.NilError(t, err)
		defer tx.Rollback()

		auction := storage.NewAuction()
		auction.ID = "04324"
		assert.NilError(t, service.AuctionInsert(tx, *auction))

		bid := storage.NewBid()
		bid.AuctionID = auction.ID
		bid.ExtraPrice = 10
		bid.State = storage.BidAccepted
		assert.NilError(t, service.BidInsert(tx, *bid))

		bid.State = storage.BidPaid
		bid.StateExtra = "vciao"
		bid.PaymentID = "3"
		bid.PaymentURL = "http"
		bid.PaymentDeadline = 4
		assert.NilError(t, service.BidUpdate(tx, *bid))

		bids, err := service.Bids(tx, auction.ID)
		assert.NilError(t, err)
		assert.Equal(t, len(bids), 1)
		assert.Equal(t, bids[0], *bid)
	})

	t.Run("TestBidFindBids", func(t *testing.T) {
		tx, err := service.DB().Begin()
		assert.NilError(t, err)
		defer tx.Rollback()

		bids := []storage.Bid{}
		bids = append(bids, *storage.NewBid())
		bids = append(bids, *storage.NewBid())
		bids = append(bids, *storage.NewBid())
		bids = append(bids, *storage.NewBid())
		result := storage.FindBids(bids, storage.BidPaid)
		assert.Equal(t, len(result), 0)
		result = storage.FindBids(bids, storage.BidAccepted)
		assert.Equal(t, len(result), 4)
		result[0].State = storage.BidPaid
		result = storage.FindBids(bids, storage.BidPaid)
		assert.Equal(t, len(result), 1)
	})
}
