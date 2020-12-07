package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/golden"
)

func TestAuctionPassCreate(t *testing.T) {
	order := storage.NewAuctionPass()
	golden.Assert(t, dump.Sdump(order), t.Name()+".golden")
}
