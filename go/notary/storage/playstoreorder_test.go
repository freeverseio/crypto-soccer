package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/golden"
)

func TestPlaystoreOrderCreate(t *testing.T) {
	order := storage.NewPlaystoreOrder()
	golden.Assert(t, dump.Sdump(order), t.Name()+".golden")
}
