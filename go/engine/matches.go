package engine

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

type Matches []Match

func Load(
	tx *sql.Tx,
	timezoneIdx uint8,
	day uint8,
) (Matches, error) {
	ms, err := storage.MatchesByTimezoneIdxAndMatchDay(tx, timezoneIdx, day)
	if err != nil {
		return nil, err
	}

	var matches Matches
	for i := range ms {
		m := ms[i]
		match := NewMatch()
		match.HomeTeam.TeamID = m.HomeTeamID
		match.VisitorTeam.TeamID = m.VisitorTeamID
		match.HomeMatchLog = m.HomeMatchLog
		match.VisitorMatchLog = m.VisitorMatchLog
		homeTeamPlayers, err := storage.PlayersByTeamId(tx, m.HomeTeamID)
		if err != nil {
			return nil, err
		}
		for _, player := range homeTeamPlayers {
			match.HomeTeam.Players[player.ShirtNumber] = NewPlayerFromSkills(player.EncodedSkills.String())
		}
		visitorTeamPlayers, err := storage.PlayersByTeamId(tx, m.VisitorTeamID)
		if err != nil {
			return nil, err
		}
		for _, player := range visitorTeamPlayers {
			match.VisitorTeam.Players[player.ShirtNumber] = NewPlayerFromSkills(player.EncodedSkills.String())
		}
		matches = append(matches, *match)
	}
	return matches, nil
}

func (b Matches) Save(tx *sql.Tx) error {
	return nil
}
