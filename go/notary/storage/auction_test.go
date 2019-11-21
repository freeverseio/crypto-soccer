package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/google/uuid"
)

func TestGetAuction(t *testing.T) {
	sto, err := storage.NewSqlite3("../../../market.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{
		UUID:       uuid.New(),
		PlayerID:   big.NewInt(2),
		CurrencyID: 2,
		Price:      big.NewInt(44),
		Rnd:        big.NewInt(2233),
		ValidUntil: big.NewInt(5),
		Signature:  "0x",
		State:      storage.AUCTION_STARTED,
		StateExtra: "pippo",
	}
	err = sto.CreateAuction(auction)
	if err != nil {
		t.Fatal(err)
	}
	result, err := sto.GetAuctions()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("Wrong number of auctions %v", len(result))
	}
	if result[0].StateExtra != "pippo" {
		t.Fatalf("Wrong state extra: %v", result[0].StateExtra)
	}
}

func TestUpdateState(t *testing.T) {
	sto, err := storage.NewSqlite3("../../../market.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{
		UUID:       uuid.New(),
		PlayerID:   big.NewInt(2),
		CurrencyID: 2,
		Price:      big.NewInt(44),
		Rnd:        big.NewInt(2233),
		ValidUntil: big.NewInt(5),
		Signature:  "0x",
		State:      storage.AUCTION_STARTED,
	}
	err = sto.CreateAuction(auction)
	if err != nil {
		t.Fatal(err)
	}
	err = sto.UpdateAuctionState(auction.UUID, storage.AUCTION_FAILED, "ddd")
	if err != nil {
		t.Fatal(err)
	}
	result, err := sto.GetAuctions()
	if err != nil {
		t.Fatal(err)
	}
	if result[0].StateExtra != "ddd" {
		t.Fatalf("Wrong state extra: %v", result[0].StateExtra)
	}
}

func TestOpenAuction(t *testing.T) {
	sto, err := storage.NewSqlite3("../../../market.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{
		UUID:       uuid.New(),
		PlayerID:   big.NewInt(2),
		CurrencyID: 2,
		Price:      big.NewInt(44),
		Rnd:        big.NewInt(2233),
		ValidUntil: big.NewInt(5),
		Signature:  "0x",
		State:      storage.AUCTION_NO_BIDS,
	}
	err = sto.CreateAuction(auction)
	if err != nil {
		t.Fatal(err)
	}
	openedAuction, err := sto.GetOpenAuctions()
	if err != nil {
		t.Fatal(err)
	}
	if len(openedAuction) != 0 {
		t.Fatalf("Wrong expected %v", len(openedAuction))
	}
	auction = storage.Auction{
		UUID:       uuid.New(),
		PlayerID:   big.NewInt(2),
		CurrencyID: 2,
		Price:      big.NewInt(44),
		Rnd:        big.NewInt(2233),
		ValidUntil: big.NewInt(5),
		Signature:  "0x",
		State:      storage.AUCTION_STARTED,
	}
	err = sto.CreateAuction(auction)
	if err != nil {
		t.Fatal(err)
	}
	openedAuction, err = sto.GetOpenAuctions()
	if err != nil {
		t.Fatal(err)
	}
	if len(openedAuction) != 1 {
		t.Fatalf("Wrong expected %v", len(openedAuction))
	}
}
