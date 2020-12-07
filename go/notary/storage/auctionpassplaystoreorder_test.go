package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/golden"
)

func TestAuctionPassPlaystoreOrderCreate(t *testing.T) {
	order := storage.NewAuctionPassPlaystoreOrder()
	golden.Assert(t, dump.Sdump(order), t.Name()+".golden")
}
