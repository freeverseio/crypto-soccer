package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"gotest.tools/assert"
)

func TestMatchEventTest(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	timezone := storage.Timezone{uint8(1)}
	timezone.Insert(tx)
	country := storage.Country{timezone.TimezoneIdx, uint32(4)}
	country.Insert(tx)
	league := storage.League{timezone.TimezoneIdx, country.CountryIdx, uint32(0)}
	league.Insert(tx)
	var team storage.Team
	team.TeamID = "10"
	team.TimezoneIdx = timezone.TimezoneIdx
	team.CountryIdx = country.CountryIdx
	team.Owner = "ciao"
	team.LeagueIdx = league.LeagueIdx
	team.Insert(tx)
	matchDayIdx := uint8(3)
	matchIdx := uint8(4)
	match := storage.Match{
		TimezoneIdx:   timezone.TimezoneIdx,
		CountryIdx:    country.CountryIdx,
		LeagueIdx:     league.LeagueIdx,
		MatchDayIdx:   matchDayIdx,
		MatchIdx:      matchIdx,
		HomeTeamID:    big.NewInt(10),
		VisitorTeamID: big.NewInt(10),
	}
	match.State = storage.MatchBegin
	assert.NilError(t, match.Insert(tx))
	player := storage.Player{}
	player.PlayerId = big.NewInt(4)
	player.TeamId = team.TeamID
	blockNumber := uint64(3)
	if err = player.Insert(tx, blockNumber); err != nil {
		t.Fatal(err)
	}
	count, err := storage.MatchEventCount(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("Expected 0")
	}

	matchEvent := storage.MatchEvent{}
	matchEvent.TimezoneIdx = int(timezone.TimezoneIdx)
	matchEvent.CountryIdx = int(country.CountryIdx)
	matchEvent.LeagueIdx = int(league.LeagueIdx)
	matchEvent.MatchDayIdx = int(matchDayIdx)
	matchEvent.MatchIdx = int(matchIdx)
	matchEvent.TeamID = team.TeamID
	matchEvent.PrimaryPlayerID.String = player.PlayerId.String()
	matchEvent.PrimaryPlayerID.Valid = true
	matchEvent.Type = storage.Attack
	if err = matchEvent.Insert(tx); err != nil {
		t.Fatal(err)
	}

	count, err = storage.MatchEventCount(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("Expected 1")
	}
}
