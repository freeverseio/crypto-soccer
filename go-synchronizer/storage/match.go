package storage

import (
	"database/sql"
	"errors"
	"math/big"

	log "github.com/sirupsen/logrus"
)

type Match struct {
	TimezoneIdx   uint8
	CountryIdx    uint32
	LeagueIdx     uint32
	MatchDayIdx   uint8
	MatchIdx      uint8
	HomeTeamID    *big.Int
	VisitorTeamID *big.Int
	HomeGoals     *uint8
	VisitorGoals  *uint8
}

func (b *Storage) MatchCreate(match Match) error {
	log.Infof("[DBMS] Create Match Day %v", match)
	_, err := b.db.Exec("INSERT INTO matches (timezone_idx, country_idx, league_idx, match_day_idx, match_idx) VALUES ($1, $2, $3, $4, $5);",
		match.TimezoneIdx,
		match.CountryIdx,
		match.LeagueIdx,
		match.MatchDayIdx,
		match.MatchIdx,
	)
	if err != nil {
		return err
	}
	return nil
}

func (b *Storage) MatchSetTeams(timezoneIdx uint8, countryIdx uint32, leagueIdx uint32, matchDayIdx uint32, matchIdx uint32, homeTeamID *big.Int, visitorTeamID *big.Int) error {
	return errors.New("porca")
}

func (b *Storage) GetMatches(timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) (*[]Match, error) {
	log.Debugf("[DBMS] Get Calendar Matches timezoneIdx %v, countryIdx %v, leagueIdx %v", timezoneIdx, countryIdx, leagueIdx)
	rows, err := b.db.Query("SELECT timezone_idx, country_idx, league_idx, match_day_idx, match_idx, home_team_id, visitor_team_id, home_goals, visitor_goals FROM matches WHERE (timezone_idx == $1 AND country_idx == $2 AND league_idx == $3);", timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var matches []Match
	for rows.Next() {
		var match Match
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
		)
		if err != nil {
			return nil, err
		}
		match.HomeTeamID, _ = new(big.Int).SetString(homeTeamID.String, 10)
		match.VisitorTeamID, _ = new(big.Int).SetString(visitorTeamID.String, 10)
		matches = append(matches, match)
	}

	return &matches, nil
}
