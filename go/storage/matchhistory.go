package storage

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

type MatchState string

const (
	MatchBegin     MatchState = "begin"
	MatchHalf      MatchState = "half"
	MatchEnd       MatchState = "end"
	MatchCancelled MatchState = "cancelled"
)

type MatchHistory struct {
	Match
	BlockNumber uint64
}

func NewMatchHistory(blockNumber uint64, match Match) *MatchHistory {
	history := MatchHistory{}
	history.Match = match
	history.BlockNumber = blockNumber
	return &history
}

func (b *Match) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] Create Match  History %v", b)
	if _, err := tx.Exec(`INSERT INTO matches 
		(block_number, timezone_idx, country_idx, league_idx, match_day_idx, 
		match_idx, home_team_id, visitor_team_id,
		seed, home_goals, visitor_goals, home_match_log, visitor_match_log,state)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		b.BlockNumber,
		b.TimezoneIdx,
		b.CountryIdx,
		b.LeagueIdx,
		b.MatchDayIdx,
		b.MatchIdx,
		b.HomeTeamID,
		b.VisitorTeamID,
		b.Seed,
		b.HomeGoals,
		b.VisitorGoals,
		b.HomeMatchLog.String(),
		b.VisitorMatchLog.String(),
		b.State,
	); err != nil {
		return err
	}
	return nil
}
