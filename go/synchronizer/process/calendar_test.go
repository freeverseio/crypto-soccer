package process_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go/storage"
)

func TestGenerateCalendarOfUnexistentLeague(t *testing.T) {
	t.Parallel()
	tx, err := universedb.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	calendar, err := process.NewCalendar(bc.Contracts)
	if err != nil {
		t.Fatal(err)
	}

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
	calendarProcessor, err := process.NewCalendar(bc.Contracts)
	if err != nil {
		t.Fatal(err)
	}

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
	calendarProcessor, err := process.NewCalendar(bc.Contracts)
	if err != nil {
		t.Fatal(err)
	}
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
	if len(matches) != int(calendarProcessor.MatchDays*calendarProcessor.MatchPerDay) {
		t.Fatalf("Wrong matches %v", len(matches))
	}
}
