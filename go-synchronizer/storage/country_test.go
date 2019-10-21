package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
)

func TestCountryCount(t *testing.T) {
	storage, err := storage.NewSqlite3("../../universe.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	count, err := storage.CountryCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 result %v", count)
	}
}

func TestCountryCreate(t *testing.T) {
	sto, err := storage.NewSqlite3("../../universe.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	timezone := uint8(4)
	sto.TimezoneCreate(storage.Timezone{timezone})
	country := storage.Country{
		TimezoneIdx: timezone,
		CountryIdx:  4,
	}
	err = sto.CountryCreate(country)
	if err != nil {
		t.Fatal(err)
	}
	count, err := sto.CountryInTimezoneCount(timezone)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 result %v", count)
	}
}

func TestGetCountry(t *testing.T) {
	sto, err := storage.NewSqlite3("../../universe.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	timezone := uint8(4)
	sto.TimezoneCreate(storage.Timezone{timezone})
	country := storage.Country{
		TimezoneIdx: timezone,
		CountryIdx:  5,
	}
	err = sto.CountryCreate(country)
	if err != nil {
		t.Fatal(err)
	}
	result, err := sto.GetCountry(country.TimezoneIdx, country.CountryIdx)
	if err != nil {
		t.Fatal(err)
	}
	if result != country {
		t.Fatalf("Expected %v got %v", country, result)
	}
}
