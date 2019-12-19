package v1_test

import (
	"testing"

	v1 "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
)

func TestCreation(t *testing.T) {
	mp, err := v1.NewMarketPay()
	if err != nil {
		t.Fatal(err)
	}
	if mp == nil {
		t.Fatal("market pay instance is nil")
	}
}

func TestCreateOrder(t *testing.T) {
	mp, err := v1.NewMarketPay()
	if err != nil {
		t.Fatal(err)
	}
	name := "pippo"
	value := "134.10"
	order, err := mp.CreateOrder(name, value)
	if err != nil {
		t.Fatal(err)
	}
	if order.Status != "DRAFT" {
		t.Fatalf("Wrong order state %v", order.Status)
	}
	if order.Amount != value {
		t.Fatalf("Expected %v recived %v", value, order.Amount)
	}
}

func TestGetOrder(t *testing.T) {
	mp, err := v1.NewMarketPay()
	if err != nil {
		t.Fatal(err)
	}
	name := "pippo"
	value := "134.10"
	order, err := mp.CreateOrder(name, value)
	if err != nil {
		t.Fatal(err)
	}
	order1, err := mp.GetOrder(order.TrusteeShortlink.Hash)
	if err != nil {
		t.Fatal(err)
	}
	if order.Name != order1.Name {
		t.Fatal("Order mistmatch")
	}
}

func TestIsPaid(t *testing.T) {
	mp, err := v1.NewMarketPay()
	if err != nil {
		t.Fatal(err)
	}
	name := "pippo"
	value := "134.10"
	order, err := mp.CreateOrder(name, value)
	if err != nil {
		t.Fatal(err)
	}
	isPaid := mp.IsPaid(*order)
	if isPaid {
		t.Fatal("Should not be paid")
	}
}
func TestDraftAndFailure(t *testing.T) {
	mp := v1.NewMockMarketPay()
	name := "pippo"
	value := "134.10"
	// create an order always sets state to DRAFT
	order, err := mp.CreateOrder(name, value)
	if err != nil {
		t.Fatal(err)
	}
	if order.Status != v1.DRAFT.String() {
		t.Fatalf("Expecting DRAFT but got %v", order.Status)
	}
	mp.SetOrderStatus(v1.REJECTED)
	if o1, err := mp.GetOrder(order.TrusteeShortlink.Hash); err != nil {
		t.Fatal(err)
	} else if o1.Status != v1.REJECTED.String() {
		t.Fatalf("expecting REJECTED, but got %v", o1.Status)
	}
	mp.SetOrderStatus(v1.FAILURE)
	if o1, err := mp.GetOrder(order.TrusteeShortlink.Hash); err != nil {
		t.Fatal(err)
	} else if o1.Status != v1.FAILURE.String() {
		t.Fatalf("expecting FAILURE, but got %v", o1.Status)
	}
}
