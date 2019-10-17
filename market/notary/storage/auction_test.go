package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/market/notary/storage"
	"github.com/google/uuid"
)

func TestGetOpenAuctions(t *testing.T) {
	sto, err := storage.NewSqlite3("../../db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	result, err := sto.GetOpenAuctions()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatalf("Expected 0 but: %v", len(result))
	}
	err = sto.CreateAuction(storage.Auction{
		UUID:       uuid.New(),
		PlayerID:   big.NewInt(1),
		Price:      big.NewInt(100),
		Rnd:        big.NewInt(4353),
		ValidUntil: big.NewInt(3),
		Signature:  "ciao",
		State:      "PAID",
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err = sto.GetOpenAuctions()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatalf("Expected 0 but: %v", len(result))
	}
	err = sto.CreateAuction(storage.Auction{
		UUID:       uuid.New(),
		PlayerID:   big.NewInt(1),
		Price:      big.NewInt(100),
		Rnd:        big.NewInt(4353),
		ValidUntil: big.NewInt(3),
		Signature:  "ciao",
		State:      "STARTED",
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err = sto.GetOpenAuctions()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("Expected 1 but: %v", len(result))
	}
}

func TestGetAuctions(t *testing.T) {
	sto, err := storage.NewSqlite3("../../db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	result, err := sto.GetAuctions()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatalf("Expected 0 got %v", len(result))
	}
	err = sto.CreateAuction(storage.Auction{
		PlayerID:   big.NewInt(1),
		Price:      big.NewInt(100),
		Rnd:        big.NewInt(4353),
		ValidUntil: big.NewInt(3),
		Signature:  "ciao",
		State:      "STARTED",
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err = sto.GetAuctions()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("Expected 1 got %v", len(result))
	}
}
