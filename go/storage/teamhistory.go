package storage

import (
	"database/sql"
	"strconv"

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
