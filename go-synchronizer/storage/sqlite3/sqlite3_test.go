package sqlite3_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage/sqlite3"
)

func TestNew(t *testing.T) {
	_, err := sqlite3.New()
	if err != nil {
		t.Fatal(err)
	}
}
