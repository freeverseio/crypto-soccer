package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

func TestCountryCount(t *testing.T) {
	err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Rollback()
	count, err := s.CountryCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 result %v", count)
	}
}

func TestCountryCreate(t *testing.T) {
	err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Rollback()
	timezone := uint8(4)
	s.TimezoneCreate(storage.Timezone{timezone})
	country := storage.Country{
		TimezoneIdx: timezone,
		CountryIdx:  4,
	}
	err = s.CountryCreate(country)
	if err != nil {
		t.Fatal(err)
	}
	count, err := s.CountryInTimezoneCount(timezone)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 result %v", count)
	}
}

func TestGetCountry(t *testing.T) {
	err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Rollback()
	timezone := uint8(4)
	s.TimezoneCreate(storage.Timezone{timezone})
	country := storage.Country{
		TimezoneIdx: timezone,
		CountryIdx:  5,
	}
	err = s.CountryCreate(country)
	if err != nil {
		t.Fatal(err)
	}
	result, err := s.GetCountry(country.TimezoneIdx, country.CountryIdx)
	if err != nil {
		t.Fatal(err)
	}
	if result != country {
		t.Fatalf("Expected %v got %v", country, result)
	}
}
