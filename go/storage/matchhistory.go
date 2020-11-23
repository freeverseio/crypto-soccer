package storage

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"

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
		seed, home_goals, visitor_goals, home_teamsumskills, visitor_teamsumskills, state, state_extra, start_epoch)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);`,
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
		b.HomeTeamSumSkills,
		b.VisitorTeamSumSkills,
		b.State,
		b.StateExtra,
		b.StartEpoch,
	); err != nil {
		return err
	}
	return nil
}

func MatchesHistoriesBulkInsert(rowsToBeInserted []*MatchHistory, tx *sql.Tx) error {
	numParams := 16
	valueStrings := make([]string, 0, len(rowsToBeInserted))
	valueArgs := make([]interface{}, 0, len(rowsToBeInserted)*numParams)
	i := 0
	for _, post := range rowsToBeInserted {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*numParams+1, i*numParams+2, i*numParams+3, i*numParams+4, i*numParams+5, i*numParams+6, i*numParams+7, i*numParams+8, i*numParams+9, i*numParams+10, i*numParams+11, i*numParams+12, i*numParams+13, i*numParams+14, i*numParams+15, i*numParams+16))
		valueArgs = append(valueArgs, post.BlockNumber)
		valueArgs = append(valueArgs, post.TimezoneIdx)
		valueArgs = append(valueArgs, post.CountryIdx)
		valueArgs = append(valueArgs, post.LeagueIdx)
		valueArgs = append(valueArgs, post.MatchDayIdx)
		valueArgs = append(valueArgs, post.MatchIdx)
		valueArgs = append(valueArgs, post.HomeTeamID.String())
		valueArgs = append(valueArgs, post.VisitorTeamID.String())
		valueArgs = append(valueArgs, hex.EncodeToString(post.Seed[:]))
		valueArgs = append(valueArgs, post.HomeGoals)
		valueArgs = append(valueArgs, post.VisitorGoals)
		valueArgs = append(valueArgs, post.HomeTeamSumSkills)
		valueArgs = append(valueArgs, post.VisitorTeamSumSkills)
		valueArgs = append(valueArgs, post.State)
		valueArgs = append(valueArgs, post.StateExtra)
		valueArgs = append(valueArgs, post.StartEpoch)
		i++
	}
	stmt := fmt.Sprintf(`INSERT INTO matches_histories (
		block_number,
		timezone_idx,
		country_idx,
		league_idx,
		match_day_idx, 
		match_idx,
		home_team_id,
		visitor_team_id,
		seed,
		home_goals,
		visitor_goals,
		home_teamsumskills,
		visitor_teamsumskills,
		state,
		state_extra,
		start_epoch
		) VALUES %s
		`, strings.Join(valueStrings, ","))
	_, err := tx.Exec(stmt, valueArgs...)
	return err
}
