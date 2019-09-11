package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/market/notary/storage"
)

func TestGetBuyOrders(t *testing.T) {
	sto, err := storage.NewSqlite3("../../db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	result, err := sto.GetBuyOrders()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatalf("Expected 0 got %v", len(result))
	}

	err = sto.CreateBuyOrder(storage.BuyOrder{
		PlayerId: big.NewInt(1),
		TeamId:   big.NewInt(2),
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err = sto.GetBuyOrders()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("Expected 1 got %v", len(result))
	}

	err = sto.DeleteBuyOrder(big.NewInt(1))
	if err != nil {
		t.Fatal(err)
	}
	result, err = sto.GetBuyOrders()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatalf("Expected 0 got %v", len(result))
	}
}
