package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
	"github.com/google/uuid"
)

func TestGetbids(t *testing.T) {
	sto, err := storage.NewSqlite3("../../../market.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	auctionUuid := uuid.New()
	result, err := sto.GetBidsOfAuction(auctionUuid)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatalf("Expected 0 got %v", len(result))
	}
	err = sto.CreateAuction(storage.Auction{
		UUID:       auctionUuid,
		PlayerID:   big.NewInt(5),
		CurrencyID: 1,
		Price:      big.NewInt(3),
		Rnd:        big.NewInt(7),
		ValidUntil: big.NewInt(8),
		Signature:  "0x0",
		State:      storage.AUCTION_STARTED,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = sto.CreateBid(storage.Bid{
		Auction: auctionUuid,
		TeamID:  big.NewInt(2),
		State:   storage.BID_ACCEPTED,
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err = sto.GetBidsOfAuction(auctionUuid)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("Expected 1 got %v", len(result))
	}
	bid := result[0]
	if bid.Is2StartAuction != false {
		t.Fatal("Expected false but true")
	}
}

// 	result, err = sto.GetBids()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(result) != 1 {
// 		t.Fatalf("Expected 1 got %v", len(result))
// 	}

// 	err = sto.DeleteBet(big.NewInt(1))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err = sto.GetBids()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(result) != 0 {
// 		t.Fatalf("Expected 0 got %v", len(result))
// 	}
// }
