package marketpay_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/marketpay"
)

func TestCreation(t *testing.T) {
	mp, err := marketpay.New()
	if err != nil {
		t.Fatal(err)
	}
	if mp == nil {
		t.Fatal("market pay instance is nil")
	}
}

func TestCreateCustomer(t *testing.T) {
	mp, err := marketpay.New()
	if err != nil {
		t.Fatal(err)
	}
	customer, err := mp.CreateCustomer("+34", "657497063")
	if err != nil {
		t.Fatal(err)
	}
	if customer.Data.ID != 817 {
		t.Fatalf("Expected %v", customer)
	}
}

func TestCreateOrder(t *testing.T) {
	mp, err := marketpay.New()
	if err != nil {
		t.Fatal(err)
	}

	sellerID, err := mp.CreateCustomer("+34", "657497063")
	if err != nil {
		t.Fatal(err)
	}
	buyerID, err := mp.CreateCustomer("+34", "659853780")
	if err != nil {
		t.Fatal(err)
	}

	name := "This is a name"
	value := "100" // $$$
	orderID, err := mp.CreateOrder(
		sellerID,
		buyerID,
		name,
		value,
	)
	if err != nil {
		t.Fatal(err)
	}
	if orderID.Data.ID == 0 {
		t.Fatalf("order wrong %v", orderID)
	}
}
