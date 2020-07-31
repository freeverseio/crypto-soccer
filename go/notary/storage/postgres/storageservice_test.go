package postgres_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/storagetest"
)

func TestStorageService(t *testing.T) {
	service := postgres.NewStorageService(db)
	storagetest.TestAuctionServiceInterface(t, service)
	storagetest.TestOfferServiceInterface(t, service)
	storagetest.TestPlaystoreOrderServiceInterface(t, service)
}
