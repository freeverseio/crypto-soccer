package storage

import (
	"database/sql"
	"errors"
	"math/big"

	log "github.com/sirupsen/logrus"
)

const BotOwner = "0x0000000000000000000000000000000000000000"

type TeamState struct {
	Owner           string
	LeagueIdx       uint32
	TeamIdxInLeague uint32
	Points          uint32
	W               uint32
	D               uint32
	L               uint32
	GoalsForward    uint32
	GoalsAgainst    uint32
	PrevPerfPoints  *big.Int
	RankingPoints   *big.Int
}

type Team struct {
	TeamID      *big.Int
	Name        string
	TimezoneIdx uint8
	CountryIdx  uint32
	State       TeamState
}

func IsBotTeam(team Team) bool {
	return team.State.Owner == BotOwner
}

func (b *Storage) TeamCreate(team Team) error {
	log.Debugf("[DBMS] Create team %v", team)
	_, err := b.db.Exec("INSERT INTO teams (team_id, timezone_idx, country_idx, owner, league_idx, team_idx_in_league, name) VALUES ($1, $2, $3, $4, $5, $6, $7);",
		team.TeamID.String(),
		team.TimezoneIdx,
		team.CountryIdx,
		team.State.Owner,
		team.State.LeagueIdx,
		team.State.TeamIdxInLeague,
		team.Name,
	)
	if err != nil {
		return err
	}

	// err = b.teamHistoryAdd(team.Id, team.State)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (b *Storage) TeamCount() (uint64, error) {
	rows, err := b.db.Query("SELECT COUNT(*) FROM teams;")
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

func (b *Storage) TeamUpdate(teamID *big.Int, teamState TeamState) error {
	log.Debugf("[DBMS] + update team state %v", teamState)
	_, err := b.db.Exec(`UPDATE teams SET 
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
						ranking_points=$11
						WHERE team_id=$12`,
		teamState.Owner,
		teamState.LeagueIdx,
		teamState.TeamIdxInLeague,
		teamState.Points,
		teamState.W,
		teamState.D,
		teamState.L,
		teamState.GoalsForward,
		teamState.GoalsAgainst,
		teamState.PrevPerfPoints.String(),
		teamState.RankingPoints.String(),
		teamID.String(),
	)
	return err
}

// func (b *Storage) teamHistoryAdd(id uint64, teamState TeamState) error {
// 	log.Infof("[DBMS] + add team history %v", teamState)
// 	_, err := b.db.Exec("INSERT INTO teams_history (teamId, blockNumber, currentLeagueId, owner, posInCurrentLeagueId, posInPrevLeagueId, prevLeagueId, inBlockIndex) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);",
// 		id,
// 		teamState.BlockNumber,
// 		teamState.CurrentLeagueId,
// 		teamState.Owner,
// 		teamState.PosInCurrentLeagueId,
// 		teamState.PosInPrevLeagueId,
// 		teamState.PrevLeagueId,
// 		teamState.InBlockIndex,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
func (b *Storage) GetTeamsInLeague(timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) ([]Team, error) {
	rows, err := b.db.Query("SELECT team_id FROM teams WHERE (timezone_idx = $1 AND country_idx = $2 AND league_idx = $3);", timezoneIdx, countryIdx, leagueIdx)
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
		team, err = b.GetTeam(teamID)

		if err != nil {
			return teams, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

func (b *Storage) GetTeamID(timezoneIdx uint8, countryIdx uint32, leagueIdx uint32, teamIdxInLeague uint32) (*big.Int, error) {
	rows, err := b.db.Query("SELECT team_id FROM teams WHERE (timezone_idx = $1 AND country_idx = $2 AND league_idx = $3 AND team_idx_in_league = $4);", timezoneIdx, countryIdx, leagueIdx, teamIdxInLeague)
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

func (b *Storage) GetTeam(teamID *big.Int) (Team, error) {
	log.Debugf("[DBMS] GetTeam of teamID %v", teamID)
	var team Team
	rows, err := b.db.Query(`SELECT 
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
	name
	FROM teams WHERE (team_id = $1);`, teamID.String())
	if err != nil {
		return team, err
	}
	defer rows.Close()
	if !rows.Next() {
		return team, errors.New("Unexistent team")
	}
	team.TeamID = teamID
	var prevPerfPoints sql.NullString
	var rankingPoints sql.NullString
	err = rows.Scan(
		&team.TimezoneIdx,
		&team.CountryIdx,
		&team.State.Owner,
		&team.State.LeagueIdx,
		&team.State.TeamIdxInLeague,
		&team.State.Points,
		&team.State.W,
		&team.State.D,
		&team.State.L,
		&team.State.GoalsForward,
		&team.State.GoalsAgainst,
		&prevPerfPoints,
		&rankingPoints,
		&team.Name,
	)
	team.State.PrevPerfPoints, _ = new(big.Int).SetString(prevPerfPoints.String, 10)
	team.State.RankingPoints, _ = new(big.Int).SetString(rankingPoints.String, 10)
	if err != nil {
		return team, err
	}
	return team, nil
}
