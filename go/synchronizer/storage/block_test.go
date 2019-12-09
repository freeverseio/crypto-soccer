package storage_test

import (
	"testing"
)

func TestGetBlockNumber(t *testing.T) {
	err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Rollback()
	number, err := s.GetBlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	if number != 0 {
		t.Fatalf("Expected 0 result %v", number)
	}
}

func TestSetBlockNumber(t *testing.T) {
	err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Rollback()
	err = s.SetBlockNumber(3)
	if err != nil {
		t.Fatal(err)
	}
	number, err := s.GetBlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	if number != 3 {
		t.Fatalf("Expected 3 result %v", number)
	}
}
