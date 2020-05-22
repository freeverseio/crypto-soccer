package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/golden"
)

func TestBidNew(t *testing.T) {
	bid := storage.NewBid()
	golden.Assert(t, dump.Sdump(bid), t.Name()+".golden")
}
