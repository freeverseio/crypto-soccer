package v1_test

import (
	"testing"

	v1 "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
)

func newTestMarket() (v1.IMarketPay, error) {
	useMock := false
	test_sandboxURL := ""
	test_sandboxPublicKey := ""
	if !useMock {
		test_sandboxURL = "https://api-sandbox.truust.io/1.0"
		test_sandboxPublicKey = "pk_stage_ZkNpNElWeEg="
	}
	factory := v1.MarketPayFactory{}
	return factory.Create(test_sandboxURL, test_sandboxPublicKey)
}

func TestCreation(t *testing.T) {
	mp, err := newTestMarket()
	if err != nil {
		t.Fatal(err)
	}
	if mp == nil {
		t.Fatal("market pay instance is nil")
	}
}

func TestCreateOrder(t *testing.T) {
	mp, err := newTestMarket()
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
	mp, err := newTestMarket()
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
	mp, err := newTestMarket()
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
