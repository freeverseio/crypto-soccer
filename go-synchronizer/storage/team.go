package storage

import (
	"errors"
	"math/big"

	log "github.com/sirupsen/logrus"
)

type TeamState struct {
	Owner string
}

type Team struct {
	TeamID      *big.Int
	TimezoneIdx uint8
	CountryIdx  uint16
	State       TeamState
}

func (b *Storage) TeamCreate(team Team) error {
	log.Debugf("[DBMS] Create team %v", team)
	_, err := b.db.Exec("INSERT INTO teams (team_id, timezone_idx, country_idx, owner) VALUES ($1, $2, $3, $4);",
		team.TeamID.String(),
		team.TimezoneIdx,
		team.CountryIdx,
		team.State.Owner,
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

// func (b *Storage) TeamStateUpdate(id uint64, teamState TeamState) error {
// 	log.Infof("[DBMS] Updating team state %v", teamState)

// 	err := b.teamUpdate(id, teamState)
// 	if err != nil {
// 		return err
// 	}
// 	err = b.teamHistoryAdd(id, teamState)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (b *Storage) TeamUpdate(teamID *big.Int, teamState TeamState) error {
	log.Debugf("[DBMS] + update team state %v", teamState)
	_, err := b.db.Exec("UPDATE teams SET owner=$1 WHERE team_id=$2",
		teamState.Owner,
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

func (b *Storage) GetTeam(teamID *big.Int) (Team, error) {
	var team Team
	rows, err := b.db.Query("SELECT timezone_idx, country_idx, owner FROM teams WHERE (team_id = $1);", teamID.String())
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
		&team.State.Owner,
	)
	if err != nil {
		return team, err
	}
	return team, nil
}
