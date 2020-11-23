package storage

import (
	"database/sql"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type MatchEventType string

const (
	Attack       MatchEventType = "attack"
	YellowCard   MatchEventType = "yellow_card"
	RedCard      MatchEventType = "red_card"
	InjurySoft   MatchEventType = "injury_soft"
	InjuryHard   MatchEventType = "injury_hard"
	Substitution MatchEventType = "substytution"
)

// MatchEvent represents a row from 'public.match_events'.
type MatchEvent struct {
	TimezoneIdx       int            `json:"timezone_idx"`        // timezone_idx
	CountryIdx        int            `json:"country_idx"`         // country_idx
	LeagueIdx         int            `json:"league_idx"`          // league_idx
	MatchDayIdx       int            `json:"match_day_idx"`       // match_day_idx
	MatchIdx          int            `json:"match_idx"`           // match_idx
	Minute            int            `json:"minute"`              // minute
	Type              MatchEventType `json:"type"`                // type
	TeamID            string         `json:"team_id"`             // team_id
	ManageToShoot     bool           `json:"manage_to_shoot"`     // manage_to_shoot
	IsGoal            bool           `json:"is_goal"`             // is_goal
	PrimaryPlayerID   sql.NullString `json:"primary_player_id"`   // primary_player_id
	SecondaryPlayerID sql.NullString `json:"secondary_player_id"` // secondary_player_id
}

func MatchEventCount(tx *sql.Tx) (uint64, error) {
	rows, err := tx.Query("SELECT COUNT(*) FROM match_events;")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	var count uint64
	rows.Scan(&count)
	return count, nil
}

func MatchEventCountByTimezoneCountryLeague(tx *sql.Tx, timezone int, countryIdx int, leagueIdx int) (uint64, error) {
	rows, err := tx.Query("SELECT COUNT(*) FROM match_events WHERE timezone_idx=$1 AND country_idx=$2 AND league_idx=$3;", timezone, countryIdx, leagueIdx)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	var count uint64
	rows.Scan(&count)
	return count, nil
}

func DeleteAllMatchEvents(tx *sql.Tx, timezone int, countryIdx int, leagueIdx int) error {
	// log.Info("[DBMS] Truncate match events table")
	_, err := tx.Exec("DELETE FROM match_events WHERE timezone_idx=$1 AND country_idx=$2 AND league_idx=$3;", timezone, countryIdx, leagueIdx)
	return err
}

func MatchEventsBulkInsert(events []*MatchEvent, tx *sql.Tx) error {
	numParams := 12
	valueStrings := make([]string, 0, len(events))
	valueArgs := make([]interface{}, 0, len(events)*numParams)
	i := 0
	for _, post := range events {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*numParams+1, i*numParams+2, i*numParams+3, i*numParams+4, i*numParams+5, i*numParams+6, i*numParams+7, i*numParams+8, i*numParams+9, i*numParams+10, i*numParams+11, i*numParams+12))
		valueArgs = append(valueArgs, post.TimezoneIdx)
		valueArgs = append(valueArgs, post.CountryIdx)
		valueArgs = append(valueArgs, post.LeagueIdx)
		valueArgs = append(valueArgs, post.MatchDayIdx)
		valueArgs = append(valueArgs, post.MatchIdx)
		valueArgs = append(valueArgs, post.Minute)
		valueArgs = append(valueArgs, post.Type)
		valueArgs = append(valueArgs, post.TeamID)
		valueArgs = append(valueArgs, post.ManageToShoot)
		valueArgs = append(valueArgs, post.IsGoal)
		valueArgs = append(valueArgs, post.PrimaryPlayerID)
		valueArgs = append(valueArgs, post.SecondaryPlayerID)
		i++
	}
	stmt := fmt.Sprintf(`INSERT INTO match_events (
		timezone_idx,
		country_idx,
		league_idx,
		match_day_idx,
		match_idx,
		minute,
		type,
		team_id,
		manage_to_shoot,
		is_goal,
		primary_player_id,
		secondary_player_id
	) VALUES %s`, strings.Join(valueStrings, ","))
	_, err := tx.Exec(stmt, valueArgs...)
	return err
}

func (b *MatchEvent) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] Insert Match Event %v", b)
	_, err := tx.Exec(`INSERT INTO match_events (
		timezone_idx,
		country_idx,
		league_idx,
		match_day_idx,
		match_idx,
		minute,
		type,
		team_id,
		manage_to_shoot,
		is_goal,
		primary_player_id,
		secondary_player_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`,
		b.TimezoneIdx,
		b.CountryIdx,
		b.LeagueIdx,
		b.MatchDayIdx,
		b.MatchIdx,
		b.Minute,
		b.Type,
		b.TeamID,
		b.ManageToShoot,
		b.IsGoal,
		b.PrimaryPlayerID,
		b.SecondaryPlayerID,
	)
	return err
}
