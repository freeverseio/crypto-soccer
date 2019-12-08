package process_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

func TestGenerateCalendarOfUnexistentLeague(t *testing.T) {
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)

	calendar, err := process.NewCalendar(bc.Contracts, universedb)
	if err != nil {
		t.Fatal(err)
	}

	timezoneIdx := uint8(1)
	countryIdx := uint32(0)
	leagueIdx := uint32(0)
	err = calendar.Generate(timezoneIdx, countryIdx, leagueIdx)
	if err == nil {
		t.Fatal("Generate calendar of unexistent league")
	}
}

func TestResetCalendar(t *testing.T) {
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)

	calendarProcessor, err := process.NewCalendar(bc.Contracts, universedb)
	if err != nil {
		t.Fatal(err)
	}

	timezoneIdx := uint8(1)
	universedb.TimezoneCreate(storage.Timezone{timezoneIdx})
	countryIdx := uint32(0)
	universedb.CountryCreate(storage.Country{timezoneIdx, countryIdx})
	leagueIdx := uint32(0)
	universedb.LeagueCreate(storage.League{timezoneIdx, countryIdx, leagueIdx})
	err = calendarProcessor.Generate(timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		t.Fatal(err)
	}
	err = calendarProcessor.Reset(timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateCalendarOfExistingLeague(t *testing.T) {
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)

	calendarProcessor, err := process.NewCalendar(bc.Contracts, universedb)
	if err != nil {
		t.Fatal(err)
	}

	timezoneIdx := uint8(1)
	universedb.TimezoneCreate(storage.Timezone{timezoneIdx})
	countryIdx := uint32(0)
	universedb.CountryCreate(storage.Country{timezoneIdx, countryIdx})
	leagueIdx := uint32(0)
	universedb.LeagueCreate(storage.League{timezoneIdx, countryIdx, leagueIdx})
	err = calendarProcessor.Generate(timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		t.Fatal(err)
	}

	matches, err := universedb.GetMatches(timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != int(calendarProcessor.MatchDays*calendarProcessor.MatchPerDay) {
		t.Fatalf("Wrong matches %v", len(matches))
	}
}
