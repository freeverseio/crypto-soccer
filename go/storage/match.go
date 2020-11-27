package storage

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"

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
	TimezoneIdx          uint8
	CountryIdx           uint32
	LeagueIdx            uint32
	MatchDayIdx          uint8
	MatchIdx             uint8
	HomeTeamID           *big.Int
	VisitorTeamID        *big.Int
	Seed                 [32]byte
	HomeGoals            uint8
	VisitorGoals         uint8
	HomeTeamSumSkills    uint32
	VisitorTeamSumSkills uint32
	State                MatchState
	StateExtra           string
	StartEpoch           int64
}

type MatchStorageService interface {
	MatchesByTimezone(timezone uint8) ([]Match, error)
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
	_, err := tx.Exec("UPDATE matches SET home_team_id = NULL, visitor_team_id = NULL, home_goals = 0, visitor_goals = 0, home_teamsumskills = 0, visitor_teamsumskills = 0, state = 'begin', state_extra = '', start_epoch = 0 WHERE (timezone_idx = $1 AND country_idx = $2 AND league_idx = $3 AND match_day_idx = $4 AND match_idx = $5);",
		timezoneIdx,
		countryIdx,
		leagueIdx,
		matchDayIdx,
		matchIdx,
	)
	return err
}

func MatchSetTeams(
	tx *sql.Tx,
	timezoneIdx uint8,
	countryIdx uint32,
	leagueIdx uint32,
	matchDayIdx uint8,
	matchIdx uint8,
	homeTeamID *big.Int,
	visitorTeamID *big.Int,
	startTime *big.Int,
) error {
	if homeTeamID == nil {
		return errors.New("nill home team id")
	}
	if visitorTeamID == nil {
		return errors.New("nill visitor team id")
	}
	_, err := tx.Exec("UPDATE matches SET home_team_id = $1, visitor_team_id = $2, start_epoch =$3 WHERE (timezone_idx = $4 AND country_idx = $5 AND league_idx = $6 AND match_day_idx = $7 AND match_idx = $8);",
		homeTeamID.String(),
		visitorTeamID.String(),
		startTime.Int64(),
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
		visitor_team_id = $2,
		home_goals = $3,
		visitor_goals = $4,
		home_teamsumskills = $5,
		visitor_teamsumskills = $6,
		state = $7,
		seed = $8,
		state_extra = $9
		WHERE (timezone_idx = $10 AND country_idx = $11 AND league_idx = $12 AND match_day_idx = $13 AND match_idx = $14);`,
		b.HomeTeamID.String(),
		b.VisitorTeamID.String(),
		b.HomeGoals,
		b.VisitorGoals,
		b.HomeTeamSumSkills,
		b.VisitorTeamSumSkills,
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
		home_teamsumskills,
		visitor_teamsumskills,
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

func MatchesByTimezoneIdxCountryIdxLeagueIdx(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) ([]Match, error) {
	log.Debugf("[DBMS] Get Calendar Matches timezoneIdx %v, countryIdx %v, leagueIdx %v", timezoneIdx, countryIdx, leagueIdx)
	rows, err := tx.Query("SELECT timezone_idx, country_idx, league_idx, match_day_idx, match_idx, home_team_id, visitor_team_id, home_goals, visitor_goals, home_teamsumskills, visitor_teamsumskills FROM matches WHERE (timezone_idx = $1 AND country_idx = $2 AND league_idx = $3);", timezoneIdx, countryIdx, leagueIdx)
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
			&match.HomeTeamSumSkills,
			&match.VisitorTeamSumSkills,
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

func MatchesStartEpochByTimezone(tx *sql.Tx, timezone uint8) ([]int64, error) {
	rows, err := tx.Query(`SELECT start_epoch FROM matches WHERE timezone_idx = $1 AND league_idx = '0' AND match_idx = '0';`, timezone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	startTimes := []int64{}
	for rows.Next() {
		var value int64
		err = rows.Scan(&value)
		if err != nil {
			return nil, err
		}
		startTimes = append(startTimes, value)
	}
	sort.Slice(startTimes, func(i, j int) bool {
		return startTimes[i] < startTimes[j]
	})
	return startTimes, nil
}

func MatchesBulkInsertUpdate(rowsToBeInserted []Match, tx *sql.Tx) error {
	numParams := 15
	var err error = nil
	maxRowsToBeInserted := int(MAX_PARAMS_IN_PG_STMT / numParams)
	x := 0
	for x < len(rowsToBeInserted) {
		newX := x + maxRowsToBeInserted
		if newX > len(rowsToBeInserted) {
			newX = len(rowsToBeInserted)
		}
		currentRowsToBeInserted := rowsToBeInserted[x:newX]
		valueStrings := make([]string, 0, len(currentRowsToBeInserted))
		valueArgs := make([]interface{}, 0, len(currentRowsToBeInserted)*numParams)
		i := 0
		for _, post := range currentRowsToBeInserted {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*numParams+1, i*numParams+2, i*numParams+3, i*numParams+4, i*numParams+5, i*numParams+6, i*numParams+7, i*numParams+8, i*numParams+9, i*numParams+10, i*numParams+11, i*numParams+12, i*numParams+13, i*numParams+14, i*numParams+15))
			valueArgs = append(valueArgs, post.TimezoneIdx)
			valueArgs = append(valueArgs, post.CountryIdx)
			valueArgs = append(valueArgs, post.LeagueIdx)
			valueArgs = append(valueArgs, post.MatchDayIdx)
			valueArgs = append(valueArgs, post.MatchIdx)
			valueArgs = append(valueArgs, post.HomeTeamID.String())
			valueArgs = append(valueArgs, post.VisitorTeamID.String())
			valueArgs = append(valueArgs, post.HomeGoals)
			valueArgs = append(valueArgs, post.VisitorGoals)
			valueArgs = append(valueArgs, post.HomeTeamSumSkills)
			valueArgs = append(valueArgs, post.VisitorTeamSumSkills)
			valueArgs = append(valueArgs, post.State)
			valueArgs = append(valueArgs, hex.EncodeToString(post.Seed[:]))
			valueArgs = append(valueArgs, post.StateExtra)
			valueArgs = append(valueArgs, post.StartEpoch)
			i++
		}
		stmt := fmt.Sprintf(`INSERT INTO matches (
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
		state,
		seed,
		state_extra,
		start_epoch
		) VALUES %s
		ON CONFLICT(timezone_idx, country_idx, league_idx, match_day_idx, match_idx) DO UPDATE SET
		home_team_id = excluded.home_team_id, 
		visitor_team_id = excluded.visitor_team_id,
		home_goals = excluded.home_goals,
		visitor_goals = excluded.visitor_goals,
		home_teamsumskills = excluded.home_teamsumskills,
		visitor_teamsumskills = excluded.visitor_teamsumskills,
		state = excluded.state,
		seed = excluded.seed,
		state_extra = excluded.state_extra
		`, strings.Join(valueStrings, ","))
		_, err = tx.Exec(stmt, valueArgs...)
		if err != nil {
			return err
		}
		x = newX
	}
	return err

}
