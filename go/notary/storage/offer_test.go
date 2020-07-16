package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/golden"
)

func TestOfferNew(t *testing.T) {
	offer := storage.NewOffer()
	golden.Assert(t, dump.Sdump(offer), t.Name()+".golden")
}
