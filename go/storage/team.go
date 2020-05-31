package storage

import (
	"database/sql"
	"errors"
	"math/big"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const BotOwner = "0x0000000000000000000000000000000000000000"

type Team struct {
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

func NewTeam() *Team {
	var team Team
	team.TeamID = "0"
	team.Tactic = "340596594427581673436941882753025"
	team.MatchLog = "0"
	team.Owner = BotOwner
	return &team
}

func (b Team) IsBot() bool {
	return b.Owner == BotOwner
}

func (b *Team) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] Create team %v", b)
	_, err := tx.Exec(`
		INSERT INTO teams (
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
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
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

func TeamCount(tx *sql.Tx) (uint64, error) {
	rows, err := tx.Query("SELECT COUNT(*) FROM teams;")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	var count uint64
	err = rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (b *Team) Update(tx *sql.Tx) error {
	log.Debugf("[DBMS] + update team id %v", b.TeamID)
	_, err := tx.Exec(`UPDATE teams SET 
						owner=$1, 
						league_idx=$2, 
						team_idx_in_league=$3,
						points=$4,
						w=$5,
						d=$6,
						l=$7,
						goals_forward=$8,
						goals_against=$9,
						prev_perf_points=$10,
						ranking_points=$11,
						training_points=$12,
						name=$13,
						tactic=$14,
						match_log=$15
						WHERE team_id=$16`,
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
		b.Name,
		b.Tactic,
		b.MatchLog,
		b.TeamID,
	)
	return err
}

func TeamsByTimezoneIdxCountryIdxLeagueIdx(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) ([]Team, error) {
	rows, err := tx.Query("SELECT team_id FROM teams WHERE (timezone_idx = $1 AND country_idx = $2 AND league_idx = $3);", timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var teamsIds []string
	for rows.Next() {
		var teamID string
		err = rows.Scan(
			&teamID,
		)
		if err != nil {
			return nil, err
		}
		teamsIds = append(teamsIds, teamID)
	}
	rows.Close()
	var teams []Team
	for i := 0; i < len(teamsIds); i++ {
		teamID := teamsIds[i]
		var team Team
		team, err = TeamByTeamId(tx, teamID)

		if err != nil {
			return teams, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

func TeamIdByTimezoneIdxCountryIdxLeagueIdx(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32, leagueIdx uint32, teamIdxInLeague uint32) (*big.Int, error) {
	rows, err := tx.Query("SELECT team_id FROM teams WHERE (timezone_idx = $1 AND country_idx = $2 AND league_idx = $3 AND team_idx_in_league = $4);", timezoneIdx, countryIdx, leagueIdx, teamIdxInLeague)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	var teamID sql.NullString
	err = rows.Scan(
		&teamID,
	)
	if err != nil {
		return nil, err
	}
	result, _ := new(big.Int).SetString(teamID.String, 10)
	return result, nil
}

func TeamSetTactic(tx *sql.Tx, teamID string, tactic string) error {
	log.Debugf("[DBMS] TeamSetTactic teamID: %v, tactic: %v", teamID, tactic)
	_, err := tx.Exec("UPDATE teams SET tactic=$1 WHERE team_id = $2;", tactic, teamID)
	return err
}

func TeamByTeamId(tx *sql.Tx, teamID string) (Team, error) {
	log.Debugf("[DBMS] TeamByTeamId of teamID %v", teamID)
	var team Team
	rows, err := tx.Query(`SELECT 
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
	FROM teams WHERE (team_id = $1);`, teamID)
	if err != nil {
		return team, err
	}
	defer rows.Close()
	if !rows.Next() {
		return team, errors.New("Unexistent team")
	}
	team.TeamID = teamID
	err = rows.Scan(
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
	)
	if err != nil {
		return team, err
	}
	return team, nil
}
