package storage

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type TeamHistory struct {
	Team
	BlockNumber uint64
}

func NewTeamHistory(blockNumber uint64, team Team) *TeamHistory {
	h := TeamHistory{}
	h.BlockNumber = blockNumber
	h.Team = team
	return &h
}

func (b TeamHistory) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] Create team history %v", b)
	_, err := tx.Exec(`
		INSERT INTO teams_histories (
			block_number,
			team_id, 
			name,
			timezone_idx, 
			country_idx, 
			owner, 
			league_idx, 
			team_idx_in_league, 
			points,
			w, d, l,
			goals_forward,
			goals_against,
			prev_perf_points,
			ranking_points,
			training_points,
			tactic,
			match_log,
			is_zombie,
			leaderboard_position
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21);`,
		b.BlockNumber,
		b.TeamID,
		b.Name,
		b.TimezoneIdx,
		b.CountryIdx,
		b.Owner,
		b.LeagueIdx,
		b.TeamIdxInLeague,
		b.Points,
		b.W,
		b.D,
		b.L,
		b.GoalsForward,
		b.GoalsAgainst,
		b.PrevPerfPoints,
		strconv.FormatUint(b.RankingPoints, 10),
		b.TrainingPoints,
		b.Tactic,
		b.MatchLog,
		b.IsZombie,
		b.LeaderboardPosition,
	)
	if err != nil {
		return err
	}

	return nil
}

func TeamHistoryByTeamId(tx *sql.Tx, teamID string) ([]TeamHistory, error) {
	log.Debugf("[DBMS] TeamHistoriesByTeamId of teamID %v", teamID)
	rows, err := tx.Query(`SELECT 
	block_number,
	team_id,
	timezone_idx,
	country_idx, 
	owner, 
	league_idx, 
	team_idx_in_league, 
	points, 
	w,d,l, 
	goals_forward, 
	goals_against, 
	prev_perf_points,
	ranking_points,
	name,
	training_points,
	tactic,
	match_log,
	is_zombie,
	leaderboard_position
	FROM teams_histories WHERE (team_id = $1);`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	histories := []TeamHistory{}
	for rows.Next() {
		team := TeamHistory{}
		if err := rows.Scan(
			&team.BlockNumber,
			&team.TeamID,
			&team.TimezoneIdx,
			&team.CountryIdx,
			&team.Owner,
			&team.LeagueIdx,
			&team.TeamIdxInLeague,
			&team.Points,
			&team.W,
			&team.D,
			&team.L,
			&team.GoalsForward,
			&team.GoalsAgainst,
			&team.PrevPerfPoints,
			&team.RankingPoints,
			&team.Name,
			&team.TrainingPoints,
			&team.Tactic,
			&team.MatchLog,
			&team.IsZombie,
			&team.LeaderboardPosition,
		); err != nil {
			return nil, err
		}
		histories = append(histories, team)
	}
	return histories, nil
}

func TeamsHistoriesBulkInsert(rowsToBeInserted []*TeamHistory, tx *sql.Tx) error {
	numParams := 21
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
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*numParams+1, i*numParams+2, i*numParams+3, i*numParams+4, i*numParams+5, i*numParams+6, i*numParams+7, i*numParams+8, i*numParams+9, i*numParams+10, i*numParams+11, i*numParams+12, i*numParams+13, i*numParams+14, i*numParams+15, i*numParams+16, i*numParams+17, i*numParams+18, i*numParams+19, i*numParams+20, i*numParams+21))
			valueArgs = append(valueArgs, post.BlockNumber)
			valueArgs = append(valueArgs, post.TeamID)
			valueArgs = append(valueArgs, post.Name)
			valueArgs = append(valueArgs, post.TimezoneIdx)
			valueArgs = append(valueArgs, post.CountryIdx)
			valueArgs = append(valueArgs, post.Owner)
			valueArgs = append(valueArgs, post.LeagueIdx)
			valueArgs = append(valueArgs, post.TeamIdxInLeague)
			valueArgs = append(valueArgs, post.Points)
			valueArgs = append(valueArgs, post.W)
			valueArgs = append(valueArgs, post.D)
			valueArgs = append(valueArgs, post.L)
			valueArgs = append(valueArgs, post.GoalsForward)
			valueArgs = append(valueArgs, post.GoalsAgainst)
			valueArgs = append(valueArgs, post.PrevPerfPoints)
			valueArgs = append(valueArgs, strconv.FormatUint(post.RankingPoints, 10))
			valueArgs = append(valueArgs, post.TrainingPoints)
			valueArgs = append(valueArgs, post.Tactic)
			valueArgs = append(valueArgs, post.MatchLog)
			valueArgs = append(valueArgs, post.IsZombie)
			valueArgs = append(valueArgs, post.LeaderboardPosition)
			i++
		}
		stmt := fmt.Sprintf(`INSERT INTO teams_histories (
			block_number,
			team_id, 
			name,
			timezone_idx, 
			country_idx, 
			owner, 
			league_idx, 
			team_idx_in_league, 
			points,
			w,
			d,
			l,
			goals_forward,
			goals_against,
			prev_perf_points,
			ranking_points,
			training_points,
			tactic,
			match_log,
			is_zombie,
			leaderboard_position
			) VALUES %s
			`, strings.Join(valueStrings, ","))
		_, err = tx.Exec(stmt, valueArgs...)
		if err != nil {
			return err
		}
		x = newX
	}
	return err
}
