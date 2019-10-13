package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/universe/synchronizer/storage"
)

func TestLeagueCount(t *testing.T) {
	storage, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	count, err := storage.LeagueCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 result %v", count)
	}
}

func TestLeagueCreate(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	timezone := uint8(1)
	countryIdx := uint32(4)
	sto.TimezoneCreate(storage.Timezone{timezone})
	sto.CountryCreate(storage.Country{timezone, countryIdx})
	league := storage.League{
		TimezoneIdx: timezone,
		CountryIdx:  countryIdx,
		LeagueIdx:   2,
	}
	err = sto.LeagueCreate(league)
	if err != nil {
		t.Fatal(err)
	}
	count, err := sto.LeagueInCountryCount(timezone, countryIdx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 result %v", count)
	}
}

// func TestGetLeague(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	league := storage.League{
// 		Id: 3,
// 	}
// 	err = sto.LeagueAdd(league)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err := sto.GetLeague(league.Id)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if result != league {
// 		t.Fatalf("Expected %v got %v", league, result)
// 	}
// }
