package v2_test

import (
	"testing"

	v2 "github.com/freeverseio/crypto-soccer/go/marketpay/v2"
)

func TestCreation(t *testing.T) {
	mp, err := v2.New()
	if err != nil {
		t.Fatal(err)
	}
	if mp == nil {
		t.Fatal("market pay instance is nil")
	}
}

func TestCreateCustomer(t *testing.T) {
	mp, err := v2.New()
	if err != nil {
		t.Fatal(err)
	}
	customer, err := mp.CreateCustomer("+34", "657497063")
	if err != nil {
		t.Fatal(err)
	}
	if customer.Data.ID == 0 {
		t.Fatal("customer Id is 0")
	}
}

func TestCreateOrder(t *testing.T) {
	mp, err := v2.New()
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
	value := "100.43" // $$$
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
	if orderID.Data.SellerLink == "" {
		t.Fatal("Seller link is empty")
	}
	if orderID.Data.BuyerLink == "" {
		t.Fatal("Buyer link is empty")
	}
}

func TestIsPaid(t *testing.T) {
	mp, err := v2.New()
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
	value := "100.43" // $$$
	order, err := mp.CreateOrder(
		sellerID,
		buyerID,
		name,
		value,
	)
	if err != nil {
		t.Fatal(err)
	}
	order, err = mp.GetOrder(order.Data.ID)
	if err != nil {
		t.Fatal(err)
	}
	isPaid := mp.IsPaid(order)
	if isPaid {
		t.Fatal("The order shouldn't be paid")
	}
}
