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
