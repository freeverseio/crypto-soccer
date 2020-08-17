package storagetest

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func testOfferServiceInterface(t *testing.T, service storage.StorageService) {
	t.Run("TestOfferByIDUnexistent", func(t *testing.T) {
		tx, err := service.DB().Begin()
		assert.NilError(t, err)
		defer tx.Rollback()
		offer, err := service.Offer(tx, "4343")
		assert.NilError(t, err)
		assert.Assert(t, offer == nil)
	})

	t.Run("TestOfferInsert", func(t *testing.T) {
		tx, err := service.DB().Begin()
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

		err = service.OfferInsert(tx, *offer)
		assert.NilError(t, err)

		result, err := service.Offer(tx, offer.ID)
		assert.NilError(t, err)
		assert.Equal(t, *result, *offer)
	})

	t.Run("TestOfferUpdate", func(t *testing.T) {
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

		auctionResult, err := service.Auction(tx, auction.ID)
		assert.NilError(t, err)
		assert.Equal(t, *auctionResult, *auction)

		offer := storage.NewOffer()
		offer.State = storage.OfferStarted
		offer.StateExtra = "priva"
		offer.Seller = "yo"
		err = service.OfferInsert(tx, *offer)
		assert.NilError(t, err)
		result, err := service.Offer(tx, offer.ID)
		assert.NilError(t, err)
		assert.Equal(t, result.State, storage.OfferStarted)
		assert.Equal(t, result.StateExtra, "priva")
		assert.Equal(t, result.Seller, "yo")

		offer.StateExtra = "privato"
		offer.Seller = "yo2"
		offer.AuctionID = "ciao"

		assert.NilError(t, service.OfferUpdate(tx, *offer))
		result, err = service.Offer(tx, offer.ID)

		assert.Equal(t, result.AuctionID, "ciao")
		assert.NilError(t, err)
		assert.Equal(t, *result, *offer)
	})

	t.Run("TestInsertSameOrderTwice", func(t *testing.T) {
		t.Skip("TODO reactive me when id is the hash")
		tx, err := service.DB().Begin()
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
		err = service.OfferInsert(tx, *offer)
		assert.NilError(t, err)
		err = service.OfferInsert(tx, *offer)
		assert.Error(t, err, "some error on duplication")
	})

	t.Run("TestPendingOffer", func(t *testing.T) {
		tx, err := service.DB().Begin()
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

		err = service.OfferInsert(tx, *offer)
		assert.NilError(t, err)

		result, err := service.Offer(tx, offer.ID)
		assert.NilError(t, err)
		assert.Equal(t, *result, *offer)

		offers, err := service.OfferPendingOffers(tx)
		assert.NilError(t, err)
		assert.Equal(t, len(offers), 2)
	})
}
