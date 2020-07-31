package postgres_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/storagetest"
)

func TestStorageHistoryService(t *testing.T) {
	service := postgres.NewStorageHistoryService(db)
	storagetest.TestStorageService(t, service)
}
