package storagetest

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func TestStorageService(t *testing.T, service storage.StorageService) {
	testAuctionServiceInterface(t, service)
	testOfferServiceInterface(t, service)
	testPlaystoreOrderServiceInterface(t, service)
}
