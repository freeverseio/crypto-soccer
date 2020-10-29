package storagetest

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func testOfferServiceInterface(t *testing.T, service storage.StorageService) {
	t.Run("TestOfferByIDUnexistent", func(t *testing.T) {
		tx, err := service.Begin()
		assert.NilError(t, err)
		defer tx.Rollback()
		offer, err := tx.Offer("4343")
		assert.Error(t, err, "Could not find the offer you queried by auctionID")
		assert.Assert(t, offer == nil)
	})

	t.Run("TestOfferInsert", func(t *testing.T) {
		tx, err := service.Begin()
		assert.NilError(t, err)
		defer tx.Rollback()
		offer := storage.NewOffer()
		offer.Rnd = 4
		offer.PlayerID = "3"
		offer.CurrencyID = 3
		offer.Price = 3
		offer.ValidUntil = 3
		offer.Signature = "3"
		offer.State = storage.OfferStarted
		offer.StateExtra = "3"
		offer.Seller = "3"
		offer.Buyer = "4"

		err = tx.OfferInsert(*offer)
		assert.Error(t, err, "Trying to insert an auction with empty auctionID")
		offer.AuctionID = "dummyAuctionID"
		err = tx.OfferInsert(*offer)
		assert.NilError(t, err)

		result, err := tx.Offer(offer.AuctionID)
		assert.NilError(t, err)
		assert.Equal(t, *result, *offer)
	})

	t.Run("TestOfferUpdate", func(t *testing.T) {
		tx, err := service.Begin()
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
		assert.NilError(t, tx.AuctionInsert(*auction))

		auctionResult, err := tx.Auction(auction.ID)
		assert.NilError(t, err)
		assert.Equal(t, *auctionResult, *auction)

		offer := storage.NewOffer()
		offer.AuctionID = "ciao"
		offer.State = storage.OfferStarted
		offer.StateExtra = "priva"
		offer.Seller = "yo"
		err = tx.OfferInsert(*offer)
		assert.NilError(t, err)
		result, err := tx.Offer(offer.AuctionID)
		assert.NilError(t, err)
		assert.Equal(t, result.State, storage.OfferStarted)
		assert.Equal(t, result.StateExtra, "priva")
		assert.Equal(t, result.Seller, "yo")

		offer.StateExtra = "privato"
		offer.Seller = "yo2"

		assert.NilError(t, tx.OfferUpdate(*offer))
		result, err = tx.Offer(offer.AuctionID)

		assert.Equal(t, result.AuctionID, "ciao")
		assert.NilError(t, err)
		assert.Equal(t, *result, *offer)
	})

	t.Run("TestInsertSameOrderTwice", func(t *testing.T) {
		tx, err := service.Begin()
		assert.NilError(t, err)
		defer tx.Rollback()

		offer := storage.NewOffer()
		offer.Rnd = 4
		offer.PlayerID = "3"
		offer.CurrencyID = 3
		offer.Price = 3
		offer.ValidUntil = 3
		offer.Signature = "3"
		offer.State = storage.OfferStarted
		offer.StateExtra = "3"
		offer.Seller = "3"
		offer.Buyer = "5"
		offer.AuctionID = "dummyAuctionID"
		err = tx.OfferInsert(*offer)
		assert.NilError(t, err)
		err = tx.OfferInsert(*offer)
		assert.Assert(t, err != nil)
	})

	t.Run("TestPendingOffer", func(t *testing.T) {
		tx, err := service.Begin()
		assert.NilError(t, err)
		defer tx.Rollback()
		offer := storage.NewOffer()
		offer.Rnd = 4
		offer.PlayerID = "3"
		offer.CurrencyID = 3
		offer.Price = 3
		offer.ValidUntil = 3
		offer.Signature = "3"
		offer.State = storage.OfferStarted
		offer.StateExtra = "3"
		offer.Seller = "3"
		offer.Buyer = "4"
		offer.AuctionID = "dummyAuctionID"

		offers, err := tx.OfferPendingOffers()
		assert.NilError(t, err)
		assert.Equal(t, len(offers), 0)

		err = tx.OfferInsert(*offer)
		assert.NilError(t, err)

		result, err := tx.Offer(offer.AuctionID)
		assert.NilError(t, err)
		assert.Equal(t, *result, *offer)

		offers, err = tx.OfferPendingOffers()
		assert.NilError(t, err)
		assert.Equal(t, len(offers), 1)
	})

	t.Run("TestOfferBy", func(t *testing.T) {
		tx, err := service.Begin()
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

		assert.NilError(t, tx.AuctionInsert(*auction))

		auctionResult, err := tx.Auction(auction.ID)
		assert.NilError(t, err)
		assert.Equal(t, *auctionResult, *auction)

		offer := storage.NewOffer()
		offer.AuctionID = "dummyAuctionID"
		offer.State = storage.OfferStarted
		offer.StateExtra = "priva"
		offer.Seller = "yo"
		err = tx.OfferInsert(*offer)
		assert.NilError(t, err)
		result, err := tx.Offer(offer.AuctionID)
		assert.NilError(t, err)
		assert.Equal(t, result.State, storage.OfferStarted)
		assert.Equal(t, result.StateExtra, "priva")
		assert.Equal(t, result.Seller, "yo")

		// if I change the auctionID it will not find it when doing an Update
		offer.StateExtra = "privato"
		offer.Seller = "yo2"
		offer.AuctionID = "ciaobella"
		err2 := tx.OfferUpdate(*offer)
		assert.Equal(t, err2 != nil, true)

		// ...or when quering
		_, err = tx.Offer(offer.AuctionID)
		assert.Error(t, err, "Could not find the offer you queried by auctionID")

		// Let's reset the original (inserted) auction ID and query successfully
		offer.AuctionID = "dummyAuctionID"
		assert.NilError(t, tx.OfferUpdate(*offer))

		result, err = tx.Offer(offer.AuctionID)
		assert.Equal(t, result.AuctionID, "dummyAuctionID")
		assert.NilError(t, err)
		assert.Equal(t, *result, *offer)

		result3, err := tx.OfferByRndPrice(int32(offer.Rnd), int32(offer.Price))
		assert.NilError(t, err)
		assert.Equal(t, *result3, *offer)

		offers, err := tx.OffersByPlayerId(offer.PlayerID)
		assert.NilError(t, err)
		assert.Equal(t, 1, len(offers))

		offers, err = tx.OffersStartedByPlayerId(offer.PlayerID)
		assert.NilError(t, err)
		assert.Equal(t, 1, len(offers))

		offer.State = storage.OfferCancelled
		assert.NilError(t, tx.OfferUpdate(*offer))

		offers, err = tx.OffersStartedByPlayerId(offer.PlayerID)
		assert.NilError(t, err)
		assert.Equal(t, 0, len(offers))

	})

	t.Run("Error on cancel offe not in started state", func(t *testing.T) {
		tx, err := service.Begin()
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
		assert.NilError(t, tx.AuctionInsert(*auction))

		offer := storage.NewOffer()
		offer.AuctionID = "dummyAuctionID"
		offer.State = storage.OfferAccepted
		offer.StateExtra = "priva"
		offer.Seller = "yo"
		err = tx.OfferInsert(*offer)
		assert.NilError(t, err)

		offer.State = storage.OfferCancelled
		err = tx.OfferUpdate(*offer)
		assert.Assert(t, err != nil)
	})

	t.Run("TestInsertOrderTwiceWithLowerPrice", func(t *testing.T) {
		tx, err := service.Begin()
		assert.NilError(t, err)
		defer tx.Rollback()

		offer := storage.NewOffer()
		offer.Rnd = 4
		offer.PlayerID = "3"
		offer.CurrencyID = 3
		offer.Price = 3
		offer.ValidUntil = 3
		offer.Signature = "3"
		offer.State = storage.OfferStarted
		offer.StateExtra = "3"
		offer.Seller = "3"
		offer.Buyer = "5"
		offer.AuctionID = "dummyAuctionID"
		err = tx.OfferInsert(*offer)
		assert.NilError(t, err)
		offer.Price = 2
		err = tx.OfferInsert(*offer)
		assert.Assert(t, err != nil)
		assert.Error(t, err, "pq: error: Price not the highest")
	})

	t.Run("TestInsertOrderTwiceWithHigherPrice", func(t *testing.T) {
		tx, err := service.Begin()
		assert.NilError(t, err)
		defer tx.Rollback()

		offer := storage.NewOffer()
		offer.Rnd = 4
		offer.PlayerID = "3"
		offer.CurrencyID = 3
		offer.Price = 3
		offer.ValidUntil = 3
		offer.Signature = "3"
		offer.State = storage.OfferStarted
		offer.StateExtra = "3"
		offer.Seller = "3"
		offer.Buyer = "5"
		offer.AuctionID = "dummyAuctionID"
		err = tx.OfferInsert(*offer)
		assert.NilError(t, err)
		offer.AuctionID = "AnotherDummyID"
		offer.Price = 4
		err = tx.OfferInsert(*offer)
		assert.NilError(t, err)
	})

	t.Run("TestOfferCancel", func(t *testing.T) {
		tx, err := service.Begin()
		assert.NilError(t, err)
		defer tx.Rollback()
		offer := storage.NewOffer()
		offer.Rnd = 4
		offer.PlayerID = "3"
		offer.CurrencyID = 3
		offer.Price = 3
		offer.ValidUntil = 3
		offer.Signature = "3"
		offer.State = storage.OfferStarted
		offer.StateExtra = "3"
		offer.Seller = "3"
		offer.Buyer = "4"
		offer.AuctionID = "dummyAuctionID"
		err = tx.OfferInsert(*offer)
		assert.NilError(t, err)

		err = tx.OfferCancel(offer.AuctionID)
		assert.NilError(t, err)

		result, err := tx.Offer(offer.AuctionID)
		assert.NilError(t, err)
		offer.State = storage.OfferCancelled
		assert.Equal(t, *result, *offer)
	})
	t.Run("TestOfferCacnelAllOffers", func(t *testing.T) {
		tx, err := service.Begin()
		assert.NilError(t, err)
		defer tx.Rollback()
		offer := storage.NewOffer()
		offer.Rnd = 4
		offer.PlayerID = "3"
		offer.CurrencyID = 3
		offer.Price = 3
		offer.ValidUntil = 3
		offer.Signature = "3"
		offer.State = storage.OfferStarted
		offer.StateExtra = "3"
		offer.Seller = "3"
		offer.Buyer = "4"
		offer.AuctionID = "dummyAuctionID"
		err = tx.OfferInsert(*offer)
		assert.NilError(t, err)

		offer.AuctionID = "dummyAuctionID2"
		offer.Price = 4
		err = tx.OfferInsert(*offer)
		assert.NilError(t, err)

		err = tx.CancelAllOffersByPlayerId(offer.PlayerID)
		assert.NilError(t, err)

		result, err := tx.Offer(offer.AuctionID)
		assert.NilError(t, err)
		assert.Equal(t, result.State, storage.OfferCancelled)
	})
}
