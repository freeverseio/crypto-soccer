package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/golden"
)

func TestAuctionNew(t *testing.T) {
	auction := storage.NewAuction()
	golden.Assert(t, dump.Sdump(auction), t.Name()+".golden")
}
