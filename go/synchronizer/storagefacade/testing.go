package storagefacade

import (
	"database/sql"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	"gotest.tools/assert"
)

func CreateMatch(t *testing.T, tx *sql.Tx) {
	timezoneIdx := uint8(1)
	countryIdx := uint32(0)
	leagueIdx := uint32(0)
	var team storage.Team
	team.TeamID = big.NewInt(10)
	team.TimezoneIdx = timezoneIdx
	team.CountryIdx = countryIdx
	team.Owner = "ciao"
	team.LeagueIdx = leagueIdx
	timezone := storage.Timezone{timezoneIdx}
	err := timezone.Insert(tx)
	assert.NilError(t, err)
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	err = country.Insert(tx)
	assert.NilError(t, err)
	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
	err = league.Insert(tx)
	assert.NilError(t, err)
	err = team.Insert(tx)
	assert.NilError(t, err)
	matchDayIdx := uint8(0)
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
	assert.NilError(t, err)
}
