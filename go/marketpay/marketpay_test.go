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
	err = mp.CreateCustomer("+34", "657497063")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateOrder(t *testing.T) {
	mp, err := marketpay.New()
	if err != nil {
		t.Fatal(err)
	}

}
