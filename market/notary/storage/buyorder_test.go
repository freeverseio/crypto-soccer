package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/market/notary/storage"
)

func TestGetBets(t *testing.T) {
	sto, err := storage.NewSqlite3("../../db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	result, err := sto.GetBets()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatalf("Expected 0 got %v", len(result))
	}

	err = sto.CreateBet(storage.Bet{
		TeamID: big.NewInt(2),
	})
	if err == nil {
		t.Fatal(err)
	}
}

// 	result, err = sto.GetBets()
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
// 	result, err = sto.GetBets()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(result) != 0 {
// 		t.Fatalf("Expected 0 got %v", len(result))
// 	}
// }
