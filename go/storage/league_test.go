package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
)

func TestLeagueCount(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	count, err := storage.LeagueCount(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 result %v", count)
	}
}

func TestLeagueCreate(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	timezone := storage.Timezone{uint8(1)}
	countryIdx := uint32(4)
	timezone.Insert(tx)
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	country.Insert(tx)
	league := storage.League{
		TimezoneIdx: timezone.TimezoneIdx,
		CountryIdx:  countryIdx,
		LeagueIdx:   2,
	}
	err = league.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}
	count, err := storage.LeagueCountByTimezoneIdxCountryIdx(tx, timezone.TimezoneIdx, countryIdx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 result %v", count)
	}
}
