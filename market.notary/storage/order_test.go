package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/market/notary/storage"
)

func TestGetOrders(t *testing.T) {
	sto, err := storage.NewSqlite3("../../market.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	result, err := sto.GetOrders()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatalf("Expected 0 got %v", len(result))
	}

	err = sto.CreateAuction(storage.Auction{
		PlayerID:   big.NewInt(1),
		Price:      big.NewInt(1000),
		Rnd:        big.NewInt(4353),
		ValidUntil: big.NewInt(3),
		State:      storage.AUCTION_STARTED,
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err = sto.GetOrders()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatalf("Expected 0 got %v", len(result))
	}

	err = sto.CreateBid(storage.Bid{
		TeamID: big.NewInt(3),
		State:  storage.BID_FILED,
	})
	if err != nil {
		t.Fatal(err)
	}
	// result, err = sto.GetOrders()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if len(result) != 1 {
	// 	t.Fatalf("Expected 1 got %v", len(result))
	// }
}
