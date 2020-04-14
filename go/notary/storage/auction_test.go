package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestAuctionByIDUnexistent(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auction, err := storage.AuctionByID(tx, "4343")
	assert.NilError(t, err)
	assert.Assert(t, auction == nil)
}

func TestAuctionInsert(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	auction := storage.NewAuction()
	auction.ID = "ciao"
	auction.Rnd = 4
	auction.PlayerID = "3"
	auction.CurrencyID = 3
	auction.Price = 3
	auction.ValidUntil = "3"
	auction.Signature = "3"
	auction.State = storage.AuctionStarted
	auction.StateExtra = "3"
	auction.PaymentURL = "3"
	auction.Seller = "3"
	assert.NilError(t, auction.Insert(tx))

	result, err := storage.AuctionByID(tx, auction.ID)
	assert.NilError(t, err)
	assert.Equal(t, *result, *auction)
}

// func TestGetAuction(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../../../market.db/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	auction := storage.Auction{
// 		ID:         uuid.New(),
// 		PlayerID:   big.NewInt(2),
// 		CurrencyID: 2,
// 		Price:      big.NewInt(44),
// 		Rnd:        big.NewInt(2233),
// 		ValidUntil: 5,
// 		Signature:  "0x",
// 		State:      storage.AUCTION_STARTED,
// 		StateExtra: "pippo",
// 	}
// 	err = sto.CreateAuction(auction)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err := sto.GetAuctions()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(result) != 1 {
// 		t.Fatalf("Wrong number of auctions %v", len(result))
// 	}
// 	if result[0].StateExtra != "pippo" {
// 		t.Fatalf("Wrong state extra: %v", result[0].StateExtra)
// 	}
// }

// func TestUpdateState(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../../../market.db/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	auction := storage.Auction{
// 		UUID:       uuid.New(),
// 		PlayerID:   big.NewInt(2),
// 		CurrencyID: 2,
// 		Price:      big.NewInt(44),
// 		Rnd:        big.NewInt(2233),
// 		ValidUntil: 5,
// 		Signature:  "0x",
// 		State:      storage.AUCTION_STARTED,
// 	}
// 	err = sto.CreateAuction(auction)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = sto.UpdateAuctionState(auction.UUID, storage.AUCTION_FAILED, "ddd")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err := sto.GetAuctions()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if result[0].StateExtra != "ddd" {
// 		t.Fatalf("Wrong state extra: %v", result[0].StateExtra)
// 	}
// }

// func TestOpenAuction(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../../../market.db/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	auction := storage.Auction{
// 		UUID:       uuid.New(),
// 		PlayerID:   big.NewInt(2),
// 		CurrencyID: 2,
// 		Price:      big.NewInt(44),
// 		Rnd:        big.NewInt(2233),
// 		ValidUntil: 5,
// 		Signature:  "0x",
// 		State:      storage.AUCTION_NO_BIDS,
// 	}
// 	err = sto.CreateAuction(auction)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	openedAuction, err := sto.GetOpenAuctions()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(openedAuction) != 0 {
// 		t.Fatalf("Wrong expected %v", len(openedAuction))
// 	}
// 	auction = storage.Auction{
// 		UUID:       uuid.New(),
// 		PlayerID:   big.NewInt(2),
// 		CurrencyID: 2,
// 		Price:      big.NewInt(44),
// 		Rnd:        big.NewInt(2233),
// 		ValidUntil: 5,
// 		Signature:  "0x",
// 		State:      storage.AUCTION_STARTED,
// 	}
// 	err = sto.CreateAuction(auction)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	openedAuction, err = sto.GetOpenAuctions()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(openedAuction) != 1 {
// 		t.Fatalf("Wrong expected %v", len(openedAuction))
// 	}
// }
