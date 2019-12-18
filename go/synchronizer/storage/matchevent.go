package storage

import "database/sql"

// MatchEvent represents a row from 'public.match_events'.
type MatchEvent struct {
	TimezoneIdx       int            `json:"timezone_idx"`        // timezone_idx
	CountryIdx        int            `json:"country_idx"`         // country_idx
	LeagueIdx         int            `json:"league_idx"`          // league_idx
	MatchDayIdx       int            `json:"match_day_idx"`       // match_day_idx
	MatchIdx          int            `json:"match_idx"`           // match_idx
	Minute            int            `json:"minute"`              // minute
	Type              string         `json:"type"`                // type
	TeamID            string         `json:"team_id"`             // team_id
	ManageToShoot     sql.NullBool   `json:"manage_to_shoot"`     // manage_to_shoot
	IsGoal            sql.NullBool   `json:"is_goal"`             // is_goal
	PrimaryPlayerID   string         `json:"primary_player_id"`   // primary_player_id
	SecondaryPlayerID sql.NullString `json:"secondary_player_id"` // secondary_player_id
}
