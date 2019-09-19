package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
)

func TestCountryCount(t *testing.T) {
	storage, err := storage.NewSqlite3("../sql/00_schema.sql")
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
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	timezone := uint8(4)
	sto.TimezoneCreate(storage.Timezone{timezone})
	country := storage.Country{
		TimezoneID: timezone,
		Idx:        4,
	}
	err = sto.CountryCreate(country)
	if err != nil {
		t.Fatal(err)
	}
	count, err := sto.CountryCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 result %v", count)
	}
}

func TestGetCountry(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	timezone := uint8(4)
	sto.TimezoneCreate(storage.Timezone{timezone})
	country := storage.Country{
		TimezoneID: timezone,
		Idx:        5,
	}
	err = sto.CountryCreate(country)
	if err != nil {
		t.Fatal(err)
	}
	result, err := sto.GetCountry(country.TimezoneID, country.Idx)
	if err != nil {
		t.Fatal(err)
	}
	if result != country {
		t.Fatalf("Expected %v got %v", country, result)
	}
}
