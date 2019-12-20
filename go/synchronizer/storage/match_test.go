package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

func TestSetMatchLogs(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	timezoneIdx := uint8(1)
	countryIdx := uint32(4)
	leagueIdx := uint32(0)
	var team storage.Team
	team.TeamID = big.NewInt(10)
	team.TimezoneIdx = timezoneIdx
	team.CountryIdx = countryIdx
	team.State.Owner = "ciao"
	team.State.LeagueIdx = leagueIdx
	timezone := storage.Timezone{timezoneIdx}
	timezone.Insert(tx)
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	country.Insert(tx)
	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
	league.Insert(tx)
	team.Insert(tx)
	matchDayIdx := uint8(3)
	matchIdx := uint8(4)
	match := storage.Match{
		TimezoneIdx:   timezoneIdx,
		CountryIdx:    countryIdx,
		LeagueIdx:     leagueIdx,
		MatchDayIdx:   matchDayIdx,
		MatchIdx:      matchIdx,
		HomeTeamID:    big.NewInt(10),
		VisitorTeamID: big.NewInt(10),
	}
	err = match.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}

	homeLog, visitorLog, err := storage.GetMatchLogs(tx, timezoneIdx, countryIdx, leagueIdx, matchDayIdx, matchIdx)
	if err != nil {
		t.Fatal(err)
	}
	if homeLog.String() != "0" {
		t.Fatalf("Home match log error %v", homeLog)
	}
	if visitorLog.String() != "0" {
		t.Fatalf("Visitor match log error %v", visitorLog)
	}

	err = storage.MatchSetResult(tx, timezoneIdx, countryIdx, leagueIdx, matchDayIdx, matchIdx, 10, 3, big.NewInt(4), big.NewInt(5))
	if err != nil {
		t.Fatal(err)
	}
	homeLog, visitorLog, err = storage.GetMatchLogs(tx, timezoneIdx, countryIdx, leagueIdx, matchDayIdx, matchIdx)
	if err != nil {
		t.Fatal(err)
	}
	if homeLog.String() != "4" {
		t.Fatalf("Home match log error %v", homeLog)
	}
	if visitorLog.String() != "5" {
		t.Fatalf("Visitor match log error %v", visitorLog)
	}
}

func TestMatchReset(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	timezoneIdx := uint8(1)
	countryIdx := uint32(4)
	leagueIdx := uint32(0)
	var team storage.Team
	team.TeamID = big.NewInt(10)
	team.TimezoneIdx = timezoneIdx
	team.CountryIdx = countryIdx
	team.State.Owner = "ciao"
	team.State.LeagueIdx = leagueIdx
	timezone := storage.Timezone{timezoneIdx}
	timezone.Insert(tx)
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	country.Insert(tx)
	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
	league.Insert(tx)
	team.Insert(tx)
	matchDayIdx := uint8(3)
	matchIdx := uint8(4)
	match := storage.Match{
		TimezoneIdx:   timezoneIdx,
		CountryIdx:    countryIdx,
		LeagueIdx:     leagueIdx,
		MatchDayIdx:   matchDayIdx,
		MatchIdx:      matchIdx,
		HomeTeamID:    big.NewInt(10),
		VisitorTeamID: big.NewInt(10),
	}
	err = match.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}
	err = storage.MatchSetResult(tx, timezoneIdx, countryIdx, leagueIdx, matchDayIdx, matchIdx, 10, 3, big.NewInt(4), big.NewInt(5))
	if err != nil {
		t.Fatal(err)
	}
	err = storage.MatchReset(tx, timezoneIdx, countryIdx, leagueIdx, matchDayIdx, matchIdx)
	if err != nil {
		t.Fatal(err)
	}
	homeLog, visitorLog, err := storage.GetMatchLogs(tx, timezoneIdx, countryIdx, leagueIdx, matchDayIdx, matchIdx)
	if err != nil {
		t.Fatal(err)
	}
	if homeLog.String() != "0" {
		t.Fatalf("Home match log error %v", homeLog)
	}
	if visitorLog.String() != "0" {
		t.Fatalf("Visitor match log error %v", visitorLog)
	}
}
