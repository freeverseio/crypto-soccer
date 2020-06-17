package memory_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage/memory"
	"github.com/freeverseio/crypto-soccer/go/storage/storagetest"
)

func TestStorageServiceInterface(t *testing.T) {
	service := memory.NewStorageService()
	storagetest.TestStorageService(t, service)
}
