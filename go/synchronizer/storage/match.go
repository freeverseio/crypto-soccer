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
	log.Debugf("[DBMS] Create Match Day %v", match)
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

func (b *Storage) MatchReset(timezoneIdx uint8, countryIdx uint32, leagueIdx uint32, matchDayIdx uint32, matchIdx uint32) error {
	_, err := b.db.Exec("UPDATE matches SET home_team_id = NULL, visitor_team_id = NULL, home_goals = NULL, visitor_goals = NULL WHERE (timezone_idx = $1 AND country_idx = $2 AND league_idx = $3 AND match_day_idx = $4 AND match_idx = $5);",
		timezoneIdx,
		countryIdx,
		leagueIdx,
		matchDayIdx,
		matchIdx,
	)
	return err
}

func (b *Storage) MatchSetTeams(timezoneIdx uint8, countryIdx uint32, leagueIdx uint32, matchDayIdx uint32, matchIdx uint32, homeTeamID *big.Int, visitorTeamID *big.Int) error {
	if homeTeamID == nil {
		return errors.New("nill home team id")
	}
	if visitorTeamID == nil {
		return errors.New("nill visitor team id")
	}
	_, err := b.db.Exec("UPDATE matches SET home_team_id = $1, visitor_team_id = $2 WHERE (timezone_idx = $3 AND country_idx = $4 AND league_idx = $5 AND match_day_idx = $6 AND match_idx = $7);",
		homeTeamID.String(),
		visitorTeamID.String(),
		timezoneIdx,
		countryIdx,
		leagueIdx,
		matchDayIdx,
		matchIdx,
	)
	return err
}

func (b *Storage) MatchSetResult(
	timezoneIdx uint8,
	countryIdx uint32,
	leagueIdx uint32,
	matchDayIdx uint8,
	matchIdx uint8,
	homeGoals uint8,
	visitorGoals uint8,
	homeMatchLog *big.Int,
	visitorMatchLog *big.Int,
) error {
	log.Debugf("[DBMS] Set result tz %v, c %v, l %v, matchDayIdx %v, matchIdx %v [ %v - %v ]", timezoneIdx, countryIdx, leagueIdx, matchDayIdx, matchIdx, homeGoals, visitorGoals)
	_, err := b.db.Exec("UPDATE matches SET home_goals = $1, visitor_goals = $2, home_match_log = $3, visitor_match_log = $4  WHERE (timezone_idx = $5 AND country_idx = $6 AND league_idx = $7 AND match_day_idx = $8 AND match_idx = $9);",
		homeGoals,
		visitorGoals,
		homeMatchLog.String(),
		visitorMatchLog.String(),
		timezoneIdx,
		countryIdx,
		leagueIdx,
		matchDayIdx,
		matchIdx,
	)
	return err
}

func (b *Storage) GetMatchLogs(
	timezoneIdx uint8,
	countryIdx uint32,
	leagueIdx uint32,
	matchDayIdx uint8,
	matchIdx uint8,
) (*big.Int, *big.Int, error) {
	log.Debugf("[DBMS] Get Match Logs timezoneIdx %v, countryIdx %v, leagueIdx %v, matchDayIdx %v, matchIdx %v", timezoneIdx, countryIdx, leagueIdx, matchDayIdx, matchIdx)
	rows, err := b.db.Query("SELECT home_match_log, visitor_match_log FROM matches WHERE (timezone_idx = $1 AND country_idx = $2 AND league_idx = $3 AND match_day_idx = $4 AND match_idx = $5);", timezoneIdx, countryIdx, leagueIdx, matchDayIdx, matchIdx)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil, errors.New("GetMatchLogs: Unexistent match")
	}
	var homeMatchLogString sql.NullString
	var visitorMatchLogString sql.NullString
	err = rows.Scan(
		&homeMatchLogString,
		&visitorMatchLogString,
	)
	if err != nil {
		return nil, nil, err
	}
	homeMatchLog, _ := new(big.Int).SetString(homeMatchLogString.String, 10)
	visitorMatchLog, _ := new(big.Int).SetString(visitorMatchLogString.String, 10)
	return homeMatchLog, visitorMatchLog, nil
}

func (b *Storage) GetMatchesInDay(timezoneIdx uint8, countryIdx uint32, leagueIdx uint32, matchDayIdx uint8) ([]Match, error) {
	log.Debugf("[DBMS] Get Calendar Matches timezoneIdx %v, countryIdx %v, leagueIdx %v", timezoneIdx, countryIdx, leagueIdx)
	rows, err := b.db.Query("SELECT match_idx, home_team_id, visitor_team_id, home_goals, visitor_goals FROM matches WHERE (timezone_idx = $1 AND country_idx = $2 AND league_idx = $3 AND match_day_idx = $4);", timezoneIdx, countryIdx, leagueIdx, matchDayIdx)
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
			&match.MatchIdx,
			&homeTeamID,
			&visitorTeamID,
			&match.HomeGoals,
			&match.VisitorGoals,
		)
		if err != nil {
			return nil, err
		}
		match.TimezoneIdx = timezoneIdx
		match.CountryIdx = countryIdx
		match.LeagueIdx = leagueIdx
		match.MatchDayIdx = matchDayIdx
		match.HomeTeamID, _ = new(big.Int).SetString(homeTeamID.String, 10)
		match.VisitorTeamID, _ = new(big.Int).SetString(visitorTeamID.String, 10)
		matches = append(matches, match)
	}
	return matches, nil
}

func (b *Storage) GetMatches(timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) ([]Match, error) {
	log.Debugf("[DBMS] Get Calendar Matches timezoneIdx %v, countryIdx %v, leagueIdx %v", timezoneIdx, countryIdx, leagueIdx)
	rows, err := b.db.Query("SELECT timezone_idx, country_idx, league_idx, match_day_idx, match_idx, home_team_id, visitor_team_id, home_goals, visitor_goals FROM matches WHERE (timezone_idx = $1 AND country_idx = $2 AND league_idx = $3);", timezoneIdx, countryIdx, leagueIdx)
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
	return matches, nil
}
