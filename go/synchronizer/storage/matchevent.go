package storage

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/freeverseio/crypto-soccer/go/matchevents"
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
	PrimaryPlayerID   string         `json:"primary_player_id"`   // primary_player_id
	SecondaryPlayerID sql.NullString `json:"secondary_player_id"` // secondary_player_id
}

func MarchEventTypeByMatchEvent(event int16) (MatchEventType, error) {
	switch event {
	case matchevents.EVNT_ATTACK:
		return Attack, nil
	case matchevents.EVNT_YELLOW:
		return YellowCard, nil
	case matchevents.EVNT_RED:
		return RedCard, nil
	case matchevents.EVNT_SOFT:
		return InjurySoft, nil
	case matchevents.EVNT_HARD:
		return InjuryHard, nil
	case matchevents.EVNT_SUBST:
		return Substitution, nil
	default:
		return "", fmt.Errorf("Unknown match event %v", event)
	}
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
