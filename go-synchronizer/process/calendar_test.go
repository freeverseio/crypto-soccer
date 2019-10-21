package process_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/testutils"
)

func TestGenerateCalendarOfUnexistentLeague(t *testing.T) {
	storage, err := storage.NewSqlite3("../../universe.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	ganache, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	ganache.DeployContracts(ganache.Owner)

	calendar, err := process.NewCalendar(ganache.Leagues, storage)
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
	sto, err := storage.NewSqlite3("../../universe.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	ganache, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	ganache.DeployContracts(ganache.Owner)

	calendarProcessor, err := process.NewCalendar(ganache.Leagues, sto)
	if err != nil {
		t.Fatal(err)
	}

	timezoneIdx := uint8(1)
	sto.TimezoneCreate(storage.Timezone{timezoneIdx})
	countryIdx := uint32(0)
	sto.CountryCreate(storage.Country{timezoneIdx, countryIdx})
	leagueIdx := uint32(0)
	sto.LeagueCreate(storage.League{timezoneIdx, countryIdx, leagueIdx})
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
	sto, err := storage.NewSqlite3("../../universe.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	ganache, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	ganache.DeployContracts(ganache.Owner)

	calendarProcessor, err := process.NewCalendar(ganache.Leagues, sto)
	if err != nil {
		t.Fatal(err)
	}

	timezoneIdx := uint8(1)
	sto.TimezoneCreate(storage.Timezone{timezoneIdx})
	countryIdx := uint32(0)
	sto.CountryCreate(storage.Country{timezoneIdx, countryIdx})
	leagueIdx := uint32(0)
	sto.LeagueCreate(storage.League{timezoneIdx, countryIdx, leagueIdx})
	err = calendarProcessor.Generate(timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		t.Fatal(err)
	}

	matches, err := sto.GetMatches(timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != int(calendarProcessor.MatchDays*calendarProcessor.MatchPerDay) {
		t.Fatalf("Wrong matches %v", len(matches))
	}
}

// func TestPopulateCalendarOfExistingLeague(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../../universe.db/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	ganache := testutils.NewGanache()
// 	ganache.DeployContracts(ganache.Owner)

// 	calendarProcessor, err := process.NewCalendar(ganache.Leagues, sto)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	timezoneIdx := uint8(1)
// 	sto.TimezoneCreate(storage.Timezone{timezoneIdx})
// 	countryIdx := uint32(0)
// 	sto.CountryCreate(storage.Country{timezoneIdx, countryIdx})
// 	leagueIdx := uint32(0)
// 	sto.LeagueCreate(storage.League{timezoneIdx, countryIdx, leagueIdx})
// 	err = calendarProcessor.Generate(timezoneIdx, countryIdx, leagueIdx)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = calendarProcessor.Populate(timezoneIdx, countryIdx, leagueIdx)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	matches, err := sto.GetMatches(timezoneIdx, countryIdx, leagueIdx)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	match := (*matches)[0]
// 	if match.HomeTeamID == nil {
// 		t.Fatal("Home team is nil")
// 	}
// }
