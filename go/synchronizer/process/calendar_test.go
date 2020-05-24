package process_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
)

func TestGenerateCalendarOfUnexistentLeague(t *testing.T) {
	t.Parallel()
	tx, err := universedb.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	calendar := process.NewCalendar(bc.Contracts)
	timezoneIdx := uint8(1)
	countryIdx := uint32(0)
	leagueIdx := uint32(0)
	err = calendar.Generate(tx, timezoneIdx, countryIdx, leagueIdx)
	if err == nil {
		t.Fatal("Generate calendar of unexistent league")
	}
}

func TestResetCalendar(t *testing.T) {
	t.Parallel()
	tx, err := universedb.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	calendarProcessor := process.NewCalendar(bc.Contracts)

	timezoneIdx := uint8(1)
	timezone := storage.Timezone{timezoneIdx}
	timezone.Insert(tx)
	countryIdx := uint32(0)
	country := storage.Country{timezoneIdx, countryIdx}
	country.Insert(tx)
	leagueIdx := uint32(0)
	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
	league.Insert(tx)
	err = calendarProcessor.Generate(tx, timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		t.Fatal(err)
	}
	err = calendarProcessor.Reset(tx, timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateCalendarOfExistingLeague(t *testing.T) {
	t.Parallel()
	tx, err := universedb.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	calendarProcessor := process.NewCalendar(bc.Contracts)
	timezoneIdx := uint8(1)
	timezone := storage.Timezone{timezoneIdx}
	timezone.Insert(tx)
	countryIdx := uint32(0)
	country := storage.Country{timezoneIdx, countryIdx}
	country.Insert(tx)
	leagueIdx := uint32(0)
	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
	league.Insert(tx)

	err = calendarProcessor.Generate(tx, timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		t.Fatal(err)
	}

	matches, err := storage.MatchesByTimezoneIdxCountryIdxLeagueIdx(tx, timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != int(contracts.MatchDays*contracts.MatchesPerDay) {
		t.Fatalf("Wrong matches %v", len(matches))
	}
}
