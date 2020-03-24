package storage

import (
	"database/sql"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type TeamHistory struct {
	BlockNumber     uint64
	TeamID          string
	TimezoneIdx     uint8
	CountryIdx      uint32
	Name            string
	Owner           string
	LeagueIdx       uint32
	TeamIdxInLeague uint32
	Points          uint32
	W               uint32
	D               uint32
	L               uint32
	GoalsForward    uint32
	GoalsAgainst    uint32
	PrevPerfPoints  uint64
	RankingPoints   uint64
	TrainingPoints  uint16
	Tactic          string
	MatchLog        string
}

func (b TeamHistory) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] Create team history %v", b)
	_, err := tx.Exec(`
		INSERT INTO teams_histories (
			block_number,
			team_id, 
			timezone_idx, 
			country_idx, 
			owner, 
			league_idx, 
			team_idx_in_league, 
			name,
			ranking_points,
			tactic,
			match_log
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`,
		b.BlockNumber,
		b.TeamID,
		b.TimezoneIdx,
		b.CountryIdx,
		b.Owner,
		b.LeagueIdx,
		b.TeamIdxInLeague,
		b.Name,
		strconv.FormatUint(b.RankingPoints, 10),
		b.Tactic,
		b.MatchLog,
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
	match_log
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
		); err != nil {
			return nil, err
		}
		histories = append(histories, team)
	}
	return histories, nil
}
