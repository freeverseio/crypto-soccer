package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/golden"
)

func TestUnpaymentCreate(t *testing.T) {
	unpayment := storage.NewUnpayment()
	golden.Assert(t, dump.Sdump(unpayment), t.Name()+".golden")
}
