package storage

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type TeamState struct {
	BlockNumber          uint64
	Owner                string
	CurrentLeagueId      uint64
	PosInCurrentLeagueId uint64
	PrevLeagueId         uint64
	PosInPrevLeagueId    uint64
}

type Team struct {
	Id                uint64
	Name              string
	CreationTimestamp string
	State             TeamState
}

func (b *Storage) TeamStateAdd(id uint64, teamState TeamState) error {
	log.Infof("(DBMS) Adding team state %v", teamState)

	err := b.teamHistoryAdd(id, teamState)
	if err != nil {
		return err
	}

	_, err = b.db.Exec("UPDATE teams SET blockNumber=$1, currentLeagueId=$2, owner=$3, posInCurrentLeagueId=$4, posInPrevLeagueId=$5, prevLeagueId=$6 WHERE id=$7;",
		teamState.BlockNumber,
		teamState.CurrentLeagueId,
		teamState.Owner,
		teamState.PosInCurrentLeagueId,
		teamState.PosInPrevLeagueId,
		teamState.PrevLeagueId,
		id,
	)
	return err
}

func (b *Storage) teamHistoryAdd(id uint64, teamState TeamState) error {
	log.Infof("(DBMS) Adding team history %v", teamState)
	_, err := b.db.Exec("INSERT INTO teams_history (teamId, blockNumber, currentLeagueId, owner, posInCurrentLeagueId, posInPrevLeagueId, prevLeagueId) VALUES ($1, $2, $3, $4, $5, $6, $7);",
		id,
		teamState.BlockNumber,
		teamState.CurrentLeagueId,
		teamState.Owner,
		teamState.PosInCurrentLeagueId,
		teamState.PosInPrevLeagueId,
		teamState.PrevLeagueId)
	if err != nil {
		return err
	}

	return nil
}

/// @TODO: retrive state from teams table
func (b *Storage) GetTeamState(id uint64) (TeamState, error) {
	teamState := TeamState{}
	rows, err := b.db.Query("SELECT blockNumber, currentLeagueId, owner, posInCurrentLeagueId, posInPrevLeagueId, prevLeagueId FROM teams WHERE id = $1;", id)
	if err != nil {
		return teamState, err
	}
	defer rows.Close()
	if !rows.Next() {
		return teamState, errors.New("Unexistent team")
	}
	rows.Scan(&teamState.BlockNumber, &teamState.CurrentLeagueId, &teamState.Owner, &teamState.PosInCurrentLeagueId, &teamState.PosInPrevLeagueId, &teamState.PrevLeagueId)

	return teamState, nil
}

func (b *Storage) TeamAdd(team Team) error {
	//  TODO: check for db is initialized
	log.Infof("(DBMS) Adding team %v", team)
	_, err := b.db.Exec("INSERT INTO teams (id, name, creationTimestamp, blockNumber, currentLeagueId, owner, posInCurrentLeagueId, posInPrevLeagueId, prevLeagueId) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);",
		team.Id,
		team.Name,
		team.CreationTimestamp,
		team.State.BlockNumber,
		team.State.CurrentLeagueId,
		team.State.Owner,
		team.State.PosInCurrentLeagueId,
		team.State.PosInPrevLeagueId,
		team.State.PrevLeagueId,
	)
	if err != nil {
		return err
	}

	err = b.teamHistoryAdd(team.Id, team.State)
	if err != nil {
		return err
	}

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
	rows.Scan(&count)
	return count, nil
}

func (b *Storage) GetTeam(id uint64) (Team, error) {
	team := Team{}
	rows, err := b.db.Query("SELECT id, name, creationTimestamp, blockNumber, currentLeagueId, owner, posInCurrentLeagueId, posInPrevLeagueId, prevLeagueId FROM teams WHERE (id = $1);", id)
	if err != nil {
		return team, err
	}
	defer rows.Close()
	if !rows.Next() {
		return team, errors.New("Unexistent team")
	}
	rows.Scan(
		&team.Id,
		&team.Name,
		&team.CreationTimestamp,
		&team.State.BlockNumber,
		&team.State.CurrentLeagueId,
		&team.State.Owner,
		&team.State.PosInCurrentLeagueId,
		&team.State.PosInPrevLeagueId,
		&team.State.PrevLeagueId,
	)
	return team, nil
}
