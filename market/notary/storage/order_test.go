package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/market/notary/storage"
)

func TestGetOrders(t *testing.T) {
	sto, err := storage.NewSqlite3("../../db/00_schema.sql")
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

	err = sto.CreateSellOrder(storage.SellOrder{
		PlayerId:   big.NewInt(1),
		Price:      1000,
		Rnd:        big.NewInt(4353),
		ValidUntil: big.NewInt(3),
		TypeOfTx:   3,
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

	err = sto.CreateBuyOrder(storage.BuyOrder{1, 3})
	if err != nil {
		t.Fatal(err)
	}
	result, err = sto.GetOrders()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("Expected 1 got %v", len(result))
	}
}
