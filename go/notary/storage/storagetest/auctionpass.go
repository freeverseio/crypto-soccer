package storagetest

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func testAuctionPassServiceInterface(t *testing.T, service storage.StorageService) {
	t.Run("insert Auction Pass playstore order", func(t *testing.T) {
		tx, err := service.Begin()
		assert.NilError(t, err)
		defer tx.Rollback()
		order := storage.NewAuctionPassPlaystoreOrder()
		assert.NilError(t, tx.AuctionPassPlayStoreInsert(*order))
	})
}
