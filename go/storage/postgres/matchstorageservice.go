package postgres

import (
	"database/sql"
	"math/big"

	"github.com/freeverseio/crypto-soccer/go/storage"
)

type MatchStorageService struct {
	tx *sql.Tx
}

func NewMatchStorageService(tx *sql.Tx) *MatchStorageService {
	return &MatchStorageService{
		tx: tx,
	}
}

// func (b MatchStorageService) Insert(match storage.Match) error {
// 	if b.matches[match.TimezoneIdx] == nil {
// 		b.matches[match.TimezoneIdx] = make(map[uint32]map[uint32]storage.Match)
// 	}
// 	if b.matches[match.TimezoneIdx][match.CountryIdx] == nil {
// 		b.matches[match.TimezoneIdx][match.CountryIdx] = make(map[uint32]storage.Match)
// 	}
// 	b.matches[match.TimezoneIdx][match.CountryIdx][match.LeagueIdx] = match
// 	return nil
// }

func (b MatchStorageService) MatchesByTimezone(timezone uint8) ([]storage.Match, error) {
	rows, err := b.tx.Query(`
		SELECT 
		timezone_idx, 
		country_idx, 
		league_idx, 
		match_day_idx, 
		match_idx, 
		home_team_id, 
		visitor_team_id, 
		home_goals, 
		visitor_goals, 
		home_teamsumskills,
		visitor_teamsumskills,
		state
		FROM matches WHERE timezone_idx = $1;`,
		timezone,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var matches []storage.Match
	for rows.Next() {
		var match storage.Match
		var homeTeamID sql.NullString
		var visitorTeamID sql.NullString
		err = rows.Scan(
			&match.TimezoneIdx,
			&match.CountryIdx,
			&match.LeagueIdx,
			&match.MatchDayIdx,
			&match.MatchIdx,
			&homeTeamID,
			&visitorTeamID,
			&match.HomeGoals,
			&match.VisitorGoals,
			&match.HomeTeamSumSkills,
			&match.VisitorTeamSumSkills,
			&match.State,
		)
		if err != nil {
			return nil, err
		}
		match.HomeTeamID, _ = new(big.Int).SetString(homeTeamID.String, 10)
		match.VisitorTeamID, _ = new(big.Int).SetString(visitorTeamID.String, 10)
		matches = append(matches, match)
	}
	return matches, nil
}
