package process_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

func TestGenerateCalendarOfUnexistentLeague(t *testing.T) {
	tx, err := universedb.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)

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
	tx, err := universedb.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)

	calendarProcessor, err := process.NewCalendar(bc.Contracts)
	if err != nil {
		t.Fatal(err)
	}

	timezoneIdx := uint8(1)
	timezone := storage.Timezone{timezoneIdx}
	timezone.TimezoneCreate(tx)
	countryIdx := uint32(0)
	country := storage.Country{timezoneIdx, countryIdx}
	country.CountryCreate(tx)
	leagueIdx := uint32(0)
	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
	league.LeagueCreate(tx)
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
	tx, err := universedb.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)

	calendarProcessor, err := process.NewCalendar(bc.Contracts)
	if err != nil {
		t.Fatal(err)
	}
	timezoneIdx := uint8(1)
	timezone := storage.Timezone{timezoneIdx}
	timezone.TimezoneCreate(tx)
	countryIdx := uint32(0)
	country := storage.Country{timezoneIdx, countryIdx}
	country.CountryCreate(tx)
	leagueIdx := uint32(0)
	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
	league.LeagueCreate(tx)

	err = calendarProcessor.Generate(tx, timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		t.Fatal(err)
	}

	matches, err := storage.GetMatches(tx, timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != int(calendarProcessor.MatchDays*calendarProcessor.MatchPerDay) {
		t.Fatalf("Wrong matches %v", len(matches))
	}
}
