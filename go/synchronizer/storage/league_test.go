package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

func TestLeagueCount(t *testing.T) {
	err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Rollback()
	count, err := s.LeagueCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 result %v", count)
	}
}

func TestLeagueCreate(t *testing.T) {
	err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Rollback()
	timezone := uint8(1)
	countryIdx := uint32(4)
	s.TimezoneCreate(storage.Timezone{timezone})
	s.CountryCreate(storage.Country{timezone, countryIdx})
	league := storage.League{
		TimezoneIdx: timezone,
		CountryIdx:  countryIdx,
		LeagueIdx:   2,
	}
	err = s.LeagueCreate(league)
	if err != nil {
		t.Fatal(err)
	}
	count, err := s.LeagueInCountryCount(timezone, countryIdx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 result %v", count)
	}
}

// func TestGetLeague(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../../../universe.db/00_schema.sql")
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
