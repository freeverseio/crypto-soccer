package processor_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/market/notary/processor"
	"github.com/freeverseio/crypto-soccer/market/notary/storage"
	"github.com/freeverseio/crypto-soccer/market/notary/testutils"
)

func TestProcess(t *testing.T) {
	sto, err := storage.NewSqlite3("../../db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	ganache := testutils.NewGanache()
	owner := ganache.CreateAccountWithBalance("1000000000000000000") // 1 eth
	ganache.DeployContracts(owner)

	processor, err := processor.NewProcessor(sto, ganache.Client, ganache.Assets)
	if err != nil {
		t.Fatal(err)
	}
	processor.Process()
}
