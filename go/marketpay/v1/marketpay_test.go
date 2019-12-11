package v1_test

import (
	"testing"

	v1 "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
)

func newTestMarket() (v1.IMarketPay, error) {
	useMock := false
	factory := v1.MarketPayFactory{}
	if !useMock {
		return factory.Create(v1.MarketPayContext{})
	}
	return factory.Create(
		v1.NewMockMarketPayContext([]v1.OrderStatus{
			v1.DRAFT,
		}))
}

func newMock(states []v1.OrderStatus) (v1.IMarketPay, error) {
	factory := v1.MarketPayFactory{}
	return factory.Create(v1.NewMockMarketPayContext(states))
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
func TestDraftAndFailure(t *testing.T) {
	mp, err := newMock([]v1.OrderStatus{v1.REJECTED, v1.FAILURE})
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
	// first time we query order we get REJECTED
	if o1, err := mp.GetOrder(order.TrusteeShortlink.Hash); err != nil {
		t.Fatal(err)
	} else if o1.Status != v1.REJECTED.String() {
		t.Fatalf("expecting REJECTED, but got %v", o1.Status)
	}
	// second time we query order we get FAILURE
	if o1, err := mp.GetOrder(order.TrusteeShortlink.Hash); err != nil {
		t.Fatal(err)
	} else if o1.Status != v1.FAILURE.String() {
		t.Fatalf("expecting FAILURE, but got %v", o1.Status)
	}
}
