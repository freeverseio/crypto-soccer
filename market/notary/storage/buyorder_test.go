package storage_test

import (
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

	err = sto.CreateBuyOrder(storage.BuyOrder{1, 2})
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

	err = sto.DeleteBuyOrder(1)
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
