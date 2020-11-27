package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const BotOwner = "0x0000000000000000000000000000000000000000"

type Team struct {
	TeamID              string
	TimezoneIdx         uint8
	CountryIdx          uint32
	Name                string
	ManagerName         string
	Owner               string
	LeagueIdx           uint32
	LeaderboardPosition int
	TeamIdxInLeague     uint32
	Points              uint32
	W                   uint32
	D                   uint32
	L                   uint32
	GoalsForward        uint32
	GoalsAgainst        uint32
	PrevPerfPoints      uint64
	RankingPoints       uint64
	TrainingPoints      uint16
	Tactic              string
	MatchLog            string
	IsZombie            bool
}

type TeamStorageService interface {
	Team(teamId string) (*Team, error)
	Insert(team Team) error
	UpdateName(teamId string, name string) error
	UpdateManagerName(teamId string, name string) error
	UpdateLeaderboardPosition(teamId string, position int) error
	TeamsByTimezoneIdxCountryIdxLeagueIdx(timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) ([]Team, error)
	TeamUpdateZombies(timezoneIdx uint8, countryIdx uint32) error
	TeamCleanZombies(timezoneIdx uint8, countryIdx uint32) error
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
			match_log,
			manager_name,
			leaderboard_position
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`,
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
		b.ManagerName,
		b.LeaderboardPosition,
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
						match_log=$15,
						manager_name=$16,
						leaderboard_position=$17
						WHERE team_id=$18`,
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
		b.ManagerName,
		b.LeaderboardPosition,
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

func TeamIdsByTimezoneIdxCountryIdxLeagueIdx(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) ([]string, error) {
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
	return teamsIds, nil
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
	match_log,
	manager_name,
	leaderboard_position,
	is_zombie
	FROM teams WHERE (team_id = $1);`, teamID)
	if err != nil {
		return team, err
	}
	defer rows.Close()
	if !rows.Next() {
		return team, errors.New("unexistent team")
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
		&team.ManagerName,
		&team.LeaderboardPosition,
		&team.IsZombie,
	)
	if err != nil {
		return team, err
	}
	return team, nil
}

func TeamUpdateZombies(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32) error {
	log.Debugf("[DBMS] TeamUpdateZombies")
	query := `
	UPDATE teams 
	SET    is_zombie = true 
	WHERE  timezone_idx = $1 AND country_idx = $2 AND team_id IN (SELECT sq.team_id 
					   FROM   (SELECT tc.team_id, 
									  Count(tc.team_id) AS 
									  number_of_tired_players_in_starting_eleven 
							   FROM   tactics tc 
									  LEFT JOIN teams t 
											 ON tc.team_id = t.team_id 
									  LEFT JOIN players p 
											 ON p.team_id = tc.team_id 
												AND ( tc.shirt_0 = p.shirt_number 
													   OR tc.shirt_1 = 
														  p.shirt_number 
													   OR tc.shirt_2 = 
														  p.shirt_number 
													   OR tc.shirt_3 = 
														  p.shirt_number 
													   OR tc.shirt_4 = 
														  p.shirt_number 
													   OR tc.shirt_5 = 
														  p.shirt_number 
													   OR tc.shirt_6 = 
														  p.shirt_number 
													   OR tc.shirt_7 = 
														  p.shirt_number 
													   OR tc.shirt_8 = 
														  p.shirt_number 
													   OR tc.shirt_9 = 
														  p.shirt_number 
													   OR tc.shirt_10 = 
														  p.shirt_number ) 
							   WHERE 
							   		t.owner <> '0x0000000000000000000000000000000000000000' 
									AND tiredness = 7 
							   		AND p.shirt_number < 25 
							   GROUP  BY tc.team_id) AS sq 
					   WHERE  number_of_tired_players_in_starting_eleven >= 9) `
	_, err := tx.Exec(query, timezoneIdx, countryIdx)
	return err
}

func TeamCleanZombies(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32) error {
	log.Debugf("[DBMS] TeamCleanZombies")
	_, err := tx.Exec("UPDATE teams SET is_zombie=false WHERE timezone_idx = $1 AND country_idx = $2 AND is_zombie=true;", timezoneIdx, countryIdx)
	return err
}

func TeamsBulkInsertUpdate(rowsToBeInserted []Team, tx *sql.Tx) error {
	numParams := 20
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
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*numParams+1, i*numParams+2, i*numParams+3, i*numParams+4, i*numParams+5, i*numParams+6, i*numParams+7, i*numParams+8, i*numParams+9, i*numParams+10, i*numParams+11, i*numParams+12, i*numParams+13, i*numParams+14, i*numParams+15, i*numParams+16, i*numParams+17, i*numParams+18, i*numParams+19, i*numParams+20))
			valueArgs = append(valueArgs, post.TeamID)
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
			valueArgs = append(valueArgs, post.Name)
			valueArgs = append(valueArgs, post.Tactic)
			valueArgs = append(valueArgs, post.MatchLog)
			valueArgs = append(valueArgs, post.ManagerName)
			valueArgs = append(valueArgs, post.LeaderboardPosition)
			valueArgs = append(valueArgs, post.TimezoneIdx)
			valueArgs = append(valueArgs, post.CountryIdx)
			i++
		}
		stmt := fmt.Sprintf(`INSERT INTO teams (
		team_id,
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
		name,
		tactic,
		match_log,
		manager_name,
		leaderboard_position,
		timezone_idx,
		country_idx
		) VALUES %s
		ON CONFLICT(team_id) DO UPDATE SET
		owner = excluded.owner, 
		league_idx = excluded.league_idx, 
		team_idx_in_league = excluded.team_idx_in_league,
		points = excluded.points,
		w = excluded.w,
		d = excluded.d,
		l = excluded.l,
		goals_forward = excluded.goals_forward,
		goals_against = excluded.goals_against,
		prev_perf_points = excluded.prev_perf_points,
		ranking_points = excluded.ranking_points,
		training_points = excluded.training_points,
		name = excluded.name,
		tactic = excluded.tactic,
		match_log = excluded.match_log,
		manager_name = excluded.manager_name,
		leaderboard_position = excluded.leaderboard_position
		`, strings.Join(valueStrings, ","))
		_, err = tx.Exec(stmt, valueArgs...)
		if err != nil {
			return err
		}
		x = newX
	}
	return err
}
