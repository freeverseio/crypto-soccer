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
	ms, err := storage.MatchesByTimezoneIdxCountryIdxLeagueIdxMatchdayIdx(tx, timezoneIdx, countryIdx, leagueIdx, day)
	if err != nil {
		return nil, err
	}

	var matches engine.Matches
	for i := range ms {
		m := ms[i]
		match := engine.NewMatch()
		match.HomeTeam.TeamID = m.HomeTeamID
		match.VisitorTeam.TeamID = m.VisitorTeamID
		match.HomeMatchLog = m.HomeMatchLog
		match.VisitorMatchLog = m.VisitorMatchLog
		homeTeamPlayers, err := storage.PlayersByTeamId(tx, m.HomeTeamID)
		if err != nil {
			return nil, err
		}
		for _, player := range homeTeamPlayers {
			match.HomeTeam.Players[player.ShirtNumber] = engine.NewPlayerFromSkills(player.EncodedSkills.String())
		}
		visitorTeamPlayers, err := storage.PlayersByTeamId(tx, m.VisitorTeamID)
		if err != nil {
			return nil, err
		}
		for _, player := range visitorTeamPlayers {
			match.VisitorTeam.Players[player.ShirtNumber] = engine.NewPlayerFromSkills(player.EncodedSkills.String())
		}
		matches = append(matches, *match)
	}

	return matches, nil
}
