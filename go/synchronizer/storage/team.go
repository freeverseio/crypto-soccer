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
	TeamID          *big.Int
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
	TrainingPoints  uint32
}

func (b *Team) Equal(team Team) bool {
	return b.TeamID.Cmp(team.TeamID) == 0 &&
		b.CountryIdx == team.CountryIdx &&
		b.TimezoneIdx == team.TimezoneIdx &&
		b.Owner == team.Owner &&
		b.Name == team.Name &&
		b.LeagueIdx == team.LeagueIdx &&
		b.TeamIdxInLeague == team.TeamIdxInLeague &&
		b.Points == team.Points &&
		b.W == team.W &&
		b.D == team.D &&
		b.L == team.L &&
		b.GoalsForward == team.GoalsForward &&
		b.GoalsAgainst == team.GoalsAgainst &&
		b.PrevPerfPoints == team.PrevPerfPoints &&
		b.RankingPoints == team.RankingPoints &&
		b.TrainingPoints == team.TrainingPoints
}

func IsBotTeam(team Team) bool {
	return team.Owner == BotOwner
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
			ranking_points
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
		b.TeamID.String(),
		b.TimezoneIdx,
		b.CountryIdx,
		b.Owner,
		b.LeagueIdx,
		b.TeamIdxInLeague,
		b.Name,
		strconv.FormatUint(b.RankingPoints, 10),
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
						name=$13
						WHERE team_id=$14`,
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
		b.TeamID.String(),
	)
	return err
}

func TeamsByTimezoneIdxCountryIdxLeagueIdx(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) ([]Team, error) {
	rows, err := tx.Query("SELECT team_id FROM teams WHERE (timezone_idx = $1 AND country_idx = $2 AND league_idx = $3);", timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var teamsIds []*big.Int
	for rows.Next() {
		var teamID sql.NullString
		err = rows.Scan(
			&teamID,
		)
		if err != nil {
			return nil, err
		}
		id, _ := new(big.Int).SetString(teamID.String, 10)
		teamsIds = append(teamsIds, id)
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

func TeamByTeamId(tx *sql.Tx, teamID *big.Int) (Team, error) {
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
	training_points
	FROM teams WHERE (team_id = $1);`, teamID.String())
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
	)
	if err != nil {
		return team, err
	}
	return team, nil
}
