package processor_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/market/go-broker/processor"
	"github.com/freeverseio/crypto-soccer/market/go-broker/storage"
)

func TestProcess(t *testing.T) {
	sto, err := storage.NewSqlite3("../../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	ethereumClient := "HTTP://127.0.0.1:8545"
	processor := processor.NewProcessor(sto, ethereumClient)
	processor.Process()
}
