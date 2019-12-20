package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

func TestCountryCount(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	count, err := storage.CountryCount(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 result %v", count)
	}
}

func TestCountryCreate(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	timezone := storage.Timezone{uint8(4)}
	timezone.Insert(tx)
	country := storage.Country{
		TimezoneIdx: timezone.TimezoneIdx,
		CountryIdx:  4,
	}
	err = country.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}
	count, err := storage.CountryInTimezoneCount(tx, timezone.TimezoneIdx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 result %v", count)
	}
}

func TestGetCountry(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	timezone := storage.Timezone{uint8(4)}
	timezone.Insert(tx)
	country := storage.Country{
		TimezoneIdx: timezone.TimezoneIdx,
		CountryIdx:  5,
	}
	err = country.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}
	result, err := storage.CountryByTimezoneIdxCountryIdx(tx, country.TimezoneIdx, country.CountryIdx)
	if err != nil {
		t.Fatal(err)
	}
	if result != country {
		t.Fatalf("Expected %v got %v", country, result)
	}
}
