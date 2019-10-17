package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/market/notary/storage"
	"github.com/google/uuid"
)

func TestUpdateAuctionState(t *testing.T) {
	sto, err := storage.NewSqlite3("../../db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{
		UUID:  uuid.New(),
		State: "PAID",
	}
	err = sto.CreateAuction(auction)
	if err != nil {
		t.Fatal(err)
	}
	newState := "STARTED"
	auction.State = newState
	err = sto.UpdateAuctionState(auction)
	if err != nil {
		t.Fatal(err)
	}
	auctions, err := sto.GetAuctions()
	if err != nil {
		t.Fatal(err)
	}
	if auctions[0].State != newState {
		t.Fatalf("Expected %v but %v", newState, auctions[0].State)
	}
}

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
		UUID:  uuid.New(),
		State: "PAID",
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
		UUID:  uuid.New(),
		State: "STARTED",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = sto.CreateAuction(storage.Auction{
		UUID:  uuid.New(),
		State: "ASSET_FROZEN",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = sto.CreateAuction(storage.Auction{
		UUID:  uuid.New(),
		State: "PAYING",
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err = sto.GetOpenAuctions()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 3 {
		t.Fatalf("Expected 3 but: %v", len(result))
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
