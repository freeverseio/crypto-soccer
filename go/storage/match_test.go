package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"gotest.tools/assert"
)

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
	team.TeamID = "10"
	team.TimezoneIdx = timezoneIdx
	team.CountryIdx = countryIdx
	team.Owner = "ciao"
	team.LeagueIdx = leagueIdx
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
	match.State = storage.MatchBegin
	err = match.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}
	err = storage.MatchSetResult(tx, timezoneIdx, countryIdx, leagueIdx, matchDayIdx, matchIdx, 10, 3)
	if err != nil {
		t.Fatal(err)
	}
	err = storage.MatchReset(tx, timezoneIdx, countryIdx, leagueIdx, matchDayIdx, matchIdx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMatchStartTimeByTimezone(t *testing.T) {
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	timezoneIdx := uint8(1)
	countryIdx := uint32(4)
	leagueIdx := uint32(0)
	var team storage.Team
	team.TeamID = "10"
	team.TimezoneIdx = timezoneIdx
	team.CountryIdx = countryIdx
	team.Owner = "ciao"
	team.LeagueIdx = leagueIdx
	timezone := storage.Timezone{timezoneIdx}
	timezone.Insert(tx)
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	country.Insert(tx)
	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
	league.Insert(tx)
	team.Insert(tx)
	matchDayIdx := uint8(0)
	match := storage.Match{
		TimezoneIdx:   timezoneIdx,
		CountryIdx:    countryIdx,
		LeagueIdx:     leagueIdx,
		MatchDayIdx:   matchDayIdx,
		HomeTeamID:    big.NewInt(10),
		VisitorTeamID: big.NewInt(10),
		State:         storage.MatchBegin,
	}
	match.MatchIdx = 0
	match.MatchDayIdx = 0
	match.StartEpoch = 44
	assert.NilError(t, match.Insert(tx))
	match.MatchIdx = 0
	match.MatchDayIdx = 1
	match.StartEpoch = 2
	assert.NilError(t, match.Insert(tx))

	startTimes, err := storage.MatchesStartEpochByTimezone(tx, timezoneIdx)
	assert.NilError(t, err)
	assert.Equal(t, len(startTimes), 2)
	assert.Equal(t, startTimes[0], int64(2))
	assert.Equal(t, startTimes[1], int64(44))
}
