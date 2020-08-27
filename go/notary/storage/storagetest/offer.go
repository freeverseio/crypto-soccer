package storagetest

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func testOfferServiceInterface(t *testing.T, service storage.StorageService) {
	t.Run("TestOfferByIDUnexistent", func(t *testing.T) {
		assert.NilError(t, service.Begin())
		defer service.Rollback()
		offer, err := service.Offer("4343")
		assert.NilError(t, err)
		assert.Assert(t, offer == nil)
	})

	t.Run("TestOfferInsert", func(t *testing.T) {
		assert.NilError(t, service.Begin())
		defer service.Rollback()
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

		err := service.OfferInsert(*offer)
		assert.NilError(t, err)

		result, err := service.Offer(offer.ID)
		assert.NilError(t, err)
		assert.Equal(t, *result, *offer)
	})

	t.Run("TestOfferUpdate", func(t *testing.T) {
		assert.NilError(t, service.Begin())
		defer service.Rollback()
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
		assert.NilError(t, service.AuctionInsert(*auction))

		auctionResult, err := service.Auction(auction.ID)
		assert.NilError(t, err)
		assert.Equal(t, *auctionResult, *auction)

		offer := storage.NewOffer()
		offer.State = storage.OfferStarted
		offer.StateExtra = "priva"
		offer.Seller = "yo"
		err = service.OfferInsert(*offer)
		assert.NilError(t, err)
		result, err := service.Offer(offer.ID)
		assert.NilError(t, err)
		assert.Equal(t, result.State, storage.OfferStarted)
		assert.Equal(t, result.StateExtra, "priva")
		assert.Equal(t, result.Seller, "yo")

		offer.StateExtra = "privato"
		offer.Seller = "yo2"
		offer.AuctionID = "ciao"

		assert.NilError(t, service.OfferUpdate(*offer))
		result, err = service.Offer(offer.ID)

		assert.Equal(t, result.AuctionID, "ciao")
		assert.NilError(t, err)
		assert.Equal(t, *result, *offer)
	})

	t.Run("TestInsertSameOrderTwice", func(t *testing.T) {
		t.Skip("TODO reactive me when id is the hash")
		assert.NilError(t, service.Begin())
		defer service.Rollback()

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
		err := service.OfferInsert(*offer)
		assert.NilError(t, err)
		err = service.OfferInsert(*offer)
		assert.Error(t, err, "some error on duplication")
	})

	t.Run("TestPendingOffer", func(t *testing.T) {
		assert.NilError(t, service.Begin())
		defer service.Rollback()
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

		offers, err := service.OfferPendingOffers()
		assert.NilError(t, err)
		assert.Equal(t, len(offers), 0)

		err = service.OfferInsert(*offer)
		assert.NilError(t, err)

		result, err := service.Offer(offer.ID)
		assert.NilError(t, err)
		assert.Equal(t, *result, *offer)

		offers, err = service.OfferPendingOffers()
		assert.NilError(t, err)
		assert.Equal(t, len(offers), 1)
	})

	t.Run("TestOfferBy", func(t *testing.T) {
		assert.NilError(t, service.Begin())
		defer service.Rollback()
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
		assert.NilError(t, service.AuctionInsert(*auction))

		auctionResult, err := service.Auction(auction.ID)
		assert.NilError(t, err)
		assert.Equal(t, *auctionResult, *auction)

		offer := storage.NewOffer()
		offer.State = storage.OfferStarted
		offer.StateExtra = "priva"
		offer.Seller = "yo"
		err = service.OfferInsert(*offer)
		assert.NilError(t, err)
		result, err := service.Offer(offer.ID)
		assert.NilError(t, err)
		assert.Equal(t, result.State, storage.OfferStarted)
		assert.Equal(t, result.StateExtra, "priva")
		assert.Equal(t, result.Seller, "yo")

		offer.StateExtra = "privato"
		offer.Seller = "yo2"
		offer.AuctionID = "ciao"

		assert.NilError(t, service.OfferUpdate(*offer))
		result, err = service.Offer(offer.ID)

		assert.Equal(t, result.AuctionID, "ciao")
		assert.NilError(t, err)
		assert.Equal(t, *result, *offer)

		result1, err := service.Offer(offer.ID)
		assert.NilError(t, err)
		assert.Equal(t, *result1, *offer)

		result2, err := service.OfferByAuctionId(offer.AuctionID)
		assert.NilError(t, err)
		assert.Equal(t, *result2, *offer)

		result3, err := service.OfferByRndPrice(int32(offer.Rnd), int32(offer.Price))
		assert.NilError(t, err)
		assert.Equal(t, *result3, *offer)

		offers, err := service.OffersByPlayerId(offer.PlayerID)
		assert.NilError(t, err)
		assert.Equal(t, 1, len(offers))

	})
}
