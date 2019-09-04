package processor_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/market/notary/processor"
	"github.com/freeverseio/crypto-soccer/market/notary/storage"
)

func TestProcess(t *testing.T) {
	sto, err := storage.NewSqlite3("../../db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	ethereumClient := "HTTP://127.0.0.1:8545"
	assetsContractAddress := "0x43e3"
	processor, err := processor.NewProcessor(sto, ethereumClient, assetsContractAddress)
	if err != nil {
		t.Fatal(err)
	}
	processor.Process()
}
