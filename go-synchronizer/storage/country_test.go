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

func TestCountryAdd(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	country := storage.Country{
		Id:          1,
		Name:        "Spain",
		TimezoneUTC: 4,
	}
	err = sto.CountryAdd(country)
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
	country := storage.Country{
		Id:          1,
		Name:        "Spain",
		TimezoneUTC: 4,
	}
	err = sto.CountryAdd(country)
	if err != nil {
		t.Fatal(err)
	}
	result, err := sto.GetCountry(country.Id)
	if err != nil {
		t.Fatal(err)
	}
	if result != country {
		t.Fatalf("Expected %v got %v", country, result)
	}
}
