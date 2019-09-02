package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/market/go-broker/storage"
)

func TestGetBuyOrders(t *testing.T) {
	sto, err := storage.NewSqlite3("../../sql/00_schema.sql")
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

	err = sto.CreateBuyOrder(storage.BuyOrder{1, 1000, "0x555"})
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
