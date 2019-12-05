package useractions_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/useractions"
)

func TestCreateScanner(t *testing.T) {
	_, err := useractions.NewScanner()
	if err != nil {
		t.Fatal(err)
	}
}
