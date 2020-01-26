package engine_test

import (
	"database/sql"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/engine"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	"gotest.tools/assert"
)

func TestNewMatchByLeagueWithNoMatches(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	timezoneIdx := uint8(1)
	day := uint8(0)
	matches, err := engine.FromStorage(tx, timezoneIdx, day)
	assert.NilError(t, err)
	assert.Equal(t, len(matches), 0)
}

func TestMatchesFromStorage(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	createMatches(t, tx)
	timezoneIdx := uint8(1)
	day := uint8(0)
	matches, err := engine.FromStorage(tx, timezoneIdx, day)
	assert.NilError(t, err)
	assert.Equal(t, len(matches), 8)
	match := matches[0]
	assert.Equal(t, match.HomeTeam.TeamID.String(), "10")
	assert.Equal(t, match.VisitorTeam.TeamID.String(), "10")
	assert.Equal(t, match.HomeMatchLog.String(), "12")
	assert.Equal(t, match.VisitorMatchLog.String(), "13")
	assert.Equal(t, match.HomeGoals, uint8(0))
	assert.Equal(t, match.VisitorGoals, uint8(0))
	assert.Equal(t, match.HomeTeam.Players[0].Skills().String(), "123456")
	assert.Equal(t, match.VisitorTeam.Players[0].Skills().String(), "123456")
	assert.Equal(t, len(match.Events), 0)
}

func TestMatchesToStorage(t *testing.T) {

}

func createMatches(t *testing.T, tx *sql.Tx) {
	timezoneIdx := uint8(1)
	timezone := storage.Timezone{timezoneIdx}
	err := timezone.Insert(tx)
	assert.NilError(t, err)

	countryIdx := uint32(0)
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	err = country.Insert(tx)
	assert.NilError(t, err)

	leagueIdx := uint32(0)
	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
	err = league.Insert(tx)
	assert.NilError(t, err)

	var team storage.Team
	team.TeamID = big.NewInt(10)
	team.TimezoneIdx = timezoneIdx
	team.CountryIdx = countryIdx
	team.Owner = "ciao"
	team.LeagueIdx = leagueIdx
	err = team.Insert(tx)
	assert.NilError(t, err)

	var player storage.Player
	player.TeamId = team.TeamID
	player.EncodedSkills = big.NewInt(123456)
	assert.NilError(t, player.Insert(tx))

	for i := 0; i < 8; i++ {
		matchDayIdx := uint8(0)
		matchIdx := uint8(i)
		match := storage.Match{
			TimezoneIdx:     timezoneIdx,
			CountryIdx:      countryIdx,
			LeagueIdx:       leagueIdx,
			MatchDayIdx:     matchDayIdx,
			MatchIdx:        matchIdx,
			HomeTeamID:      big.NewInt(10),
			VisitorTeamID:   big.NewInt(10),
			HomeMatchLog:    big.NewInt(12),
			VisitorMatchLog: big.NewInt(13),
		}
		err = match.Insert(tx)
		assert.NilError(t, err)
	}
}
