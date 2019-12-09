package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

func TestTimezoneCount(t *testing.T) {
	err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Rollback()
	count, err := s.TimezoneCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 result %v", count)
	}
	timezone := storage.Timezone{1}
	err = s.TimezoneCreate(timezone)
	if err != nil {
		t.Fatal(err)
	}
	count, err = s.TimezoneCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 result %v", count)
	}
}
