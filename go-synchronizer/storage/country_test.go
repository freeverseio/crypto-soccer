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
	if count != 2 {
		t.Fatalf("Expected 2 result %v", count)
	}
}

func TestCountryAdd(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	country := storage.Country{
		Id:          3,
		Name:        "Georgia",
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
	if count != 3 {
		t.Fatalf("Expected 3 result %v", count)
	}
}

func TestGetCountry(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	country := storage.Country{
		Id:          3,
		Name:        "Russia",
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
