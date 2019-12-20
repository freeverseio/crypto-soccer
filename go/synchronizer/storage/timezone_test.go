package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

func TestTimezoneCount(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	count, err := storage.TimezoneCount(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 result %v", count)
	}
	timezone := storage.Timezone{1}
	err = timezone.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}
	count, err = storage.TimezoneCount(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 result %v", count)
	}
}
