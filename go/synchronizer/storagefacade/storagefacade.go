package storagefacade

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/engine"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

func NewMatchesByLeague(
	tx *sql.Tx,
	timezoneIdx uint8,
	countryIdx uint32,
	leagueIdx uint32,
	day uint8,
) (engine.Matches, error) {
	stoMatches, err := storage.MatchesByTimezoneIdxCountryIdxLeagueIdxMatchdayIdx(tx, timezoneIdx, countryIdx, leagueIdx, day)
	if err != nil {
		return nil, err
	}

	var matches engine.Matches
	for i := range stoMatches {
		stoMatch := stoMatches[i]
		match := engine.NewMatch()
		match.HomeTeam.TeamID = stoMatch.HomeTeamID
		match.VisitorTeam.TeamID = stoMatch.VisitorTeamID
		match.HomeMatchLog = stoMatch.HomeMatchLog
		match.VisitorMatchLog = stoMatch.VisitorMatchLog
		matches = append(matches, *match)
	}

	return matches, nil
}
