package storage

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"math/big"

	log "github.com/sirupsen/logrus"
)

type MatchState string

const (
	MatchBegin     MatchState = "begin"
	MatchHalf      MatchState = "half"
	MatchEnd       MatchState = "end"
	MatchCancelled MatchState = "cancelled"
)

type Match struct {
	TimezoneIdx   uint8
	CountryIdx    uint32
	LeagueIdx     uint32
	MatchDayIdx   uint8
	MatchIdx      uint8
	HomeTeamID    *big.Int
	VisitorTeamID *big.Int
	Seed          [32]byte
	HomeGoals     uint8
	VisitorGoals  uint8
	State         MatchState
	StateExtra    string
	StartEpoch    int64
}

func NewMatch() *Match {
	return &Match{
		HomeTeamID:    big.NewInt(0),
		VisitorTeamID: big.NewInt(0),
		State:         MatchBegin,
	}
}

func (b *Match) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] Create Match Day %v", b)
	_, err := tx.Exec("INSERT INTO matches (timezone_idx, country_idx, league_idx, match_day_idx, match_idx, state, state_extra, start_epoch) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);",
		b.TimezoneIdx,
		b.CountryIdx,
		b.LeagueIdx,
		b.MatchDayIdx,
		b.MatchIdx,
		b.State,
		b.StateExtra,
		b.StartEpoch,
	)
	if err != nil {
		return err
	}
	return nil
}

func MatchReset(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32, leagueIdx uint32, matchDayIdx uint8, matchIdx uint8) error {
	_, err := tx.Exec("UPDATE matches SET home_team_id = NULL, visitor_team_id = NULL, home_goals = 0, visitor_goals = 0, state = 'begin', state_extra = '' WHERE (timezone_idx = $1 AND country_idx = $2 AND league_idx = $3 AND match_day_idx = $4 AND match_idx = $5);",
		timezoneIdx,
		countryIdx,
		leagueIdx,
		matchDayIdx,
		matchIdx,
	)
	return err
}

func MatchSetTeams(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32, leagueIdx uint32, matchDayIdx uint8, matchIdx uint8, homeTeamID *big.Int, visitorTeamID *big.Int) error {
	if homeTeamID == nil {
		return errors.New("nill home team id")
	}
	if visitorTeamID == nil {
		return errors.New("nill visitor team id")
	}
	_, err := tx.Exec("UPDATE matches SET home_team_id = $1, visitor_team_id = $2 WHERE (timezone_idx = $3 AND country_idx = $4 AND league_idx = $5 AND match_day_idx = $6 AND match_idx = $7);",
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

func (b Match) Update(tx *sql.Tx, blockNumber uint64) error {
	if _, err := tx.Exec(`
		UPDATE matches SET 
		home_team_id = $1, 
		visitor_team_id = $2 ,
		home_goals = $3,
		visitor_goals = $4,
		state = $5,
		seed = $6,
		state_extra = $7
		WHERE (timezone_idx = $8 AND country_idx = $9 AND league_idx = $10 AND match_day_idx = $11 AND match_idx = $12);`,
		b.HomeTeamID.String(),
		b.VisitorTeamID.String(),
		b.HomeGoals,
		b.VisitorGoals,
		b.State,
		hex.EncodeToString(b.Seed[:]),
		b.StateExtra,
		b.TimezoneIdx,
		b.CountryIdx,
		b.LeagueIdx,
		b.MatchDayIdx,
		b.MatchIdx,
	); err != nil {
		return err
	}
	history := NewMatchHistory(blockNumber, b)
	if err := history.Insert(tx); err != nil {
		return err
	}
	return nil
}

func MatchSetResult(
	tx *sql.Tx,
	timezoneIdx uint8,
	countryIdx uint32,
	leagueIdx uint32,
	matchDayIdx uint8,
	matchIdx uint8,
	homeGoals uint8,
	visitorGoals uint8,
) error {
	log.Debugf("[DBMS] Set result tz %v, c %v, l %v, matchDayIdx %v, matchIdx %v [ %v - %v ]", timezoneIdx, countryIdx, leagueIdx, matchDayIdx, matchIdx, homeGoals, visitorGoals)
	_, err := tx.Exec("UPDATE matches SET home_goals = $1, visitor_goals = $2 WHERE (timezone_idx = $3 AND country_idx = $4 AND league_idx = $5 AND match_day_idx = $6 AND match_idx = $7);",
		homeGoals,
		visitorGoals,
		timezoneIdx,
		countryIdx,
		leagueIdx,
		matchDayIdx,
		matchIdx,
	)
	return err
}

func MatchesByTimezoneIdxCountryIdxLeagueIdxMatchdayIdx(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32, leagueIdx uint32, matchDayIdx uint8) ([]Match, error) {
	log.Debugf("[DBMS] Get Calendar Matches timezoneIdx %v, countryIdx %v, leagueIdx %v", timezoneIdx, countryIdx, leagueIdx)
	var matchesInDay []Match
	matches, err := MatchesByTimezoneIdxCountryIdxLeagueIdx(tx, timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		return matchesInDay, err
	}
	for i := range matches {
		if matches[i].MatchDayIdx == matchDayIdx {
			matchesInDay = append(matchesInDay, matches[i])
		}
	}
	return matchesInDay, nil
}

func MatchesByTimezoneIdxAndMatchDay(tx *sql.Tx, timezoneIdx uint8, matchDayIdx uint8) ([]Match, error) {
	log.Debugf("[DBMS] Get Matches timezoneIdx %v, matchDayIdx %v", timezoneIdx, matchDayIdx)
	rows, err := tx.Query(`
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
		state
		FROM matches WHERE (timezone_idx = $1 AND match_day_idx = $2);`,
		timezoneIdx,
		matchDayIdx,
	)
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

func MatchesByTimezoneIdxCountryIdxLeagueIdx(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) ([]Match, error) {
	log.Debugf("[DBMS] Get Calendar Matches timezoneIdx %v, countryIdx %v, leagueIdx %v", timezoneIdx, countryIdx, leagueIdx)
	rows, err := tx.Query("SELECT timezone_idx, country_idx, league_idx, match_day_idx, match_idx, home_team_id, visitor_team_id, home_goals, visitor_goals FROM matches WHERE (timezone_idx = $1 AND country_idx = $2 AND league_idx = $3);", timezoneIdx, countryIdx, leagueIdx)
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
