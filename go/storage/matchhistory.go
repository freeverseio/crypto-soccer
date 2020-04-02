package storage

import (
	"database/sql"
	"encoding/hex"

	log "github.com/sirupsen/logrus"
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

func (b MatchHistory) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] Create Match  History %v", b)
	if _, err := tx.Exec(`INSERT INTO matches_histories 
		(block_number, timezone_idx, country_idx, league_idx, match_day_idx, 
		match_idx, home_team_id, visitor_team_id,
		seed, home_goals, visitor_goals, state, state_extra, start_epoch)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);`,
		b.BlockNumber,
		b.TimezoneIdx,
		b.CountryIdx,
		b.LeagueIdx,
		b.MatchDayIdx,
		b.MatchIdx,
		b.HomeTeamID.String(),
		b.VisitorTeamID.String(),
		hex.EncodeToString(b.Seed[:]),
		b.HomeGoals,
		b.VisitorGoals,
		b.State,
		b.StateExtra,
		b.StartEpoch,
	); err != nil {
		return err
	}
	return nil
}
