package storage

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type TeamState struct {
	BlockNumber          uint64
	InBlockIndex         uint64
	Owner                string
	CurrentLeagueId      uint64
	PosInCurrentLeagueId uint64
	PrevLeagueId         uint64
	PosInPrevLeagueId    uint64
}

type Team struct {
	Id                uint64
	Name              string
	CreationTimestamp uint64
	State             TeamState
}

func (b *Storage) TeamStateUpdate(id uint64, teamState TeamState) error {
	log.Infof("[DBMS] Updating team state %v", teamState)

	err := b.teamStateUpdate(id, teamState)
	if err != nil {
		return err
	}
	err = b.teamHistoryAdd(id, teamState)
	if err != nil {
		return err
	}
	return nil
}

func (b *Storage) teamStateUpdate(id uint64, teamState TeamState) error {
	log.Infof("[DBMS] + update team state %v", teamState)

	_, err := b.db.Exec("UPDATE teams SET blockNumber=$1, currentLeagueId=$2, owner=$3, posInCurrentLeagueId=$4, posInPrevLeagueId=$5, prevLeagueId=$6, inBlockIndex=$7 WHERE id=$8;",
		teamState.BlockNumber,
		teamState.CurrentLeagueId,
		teamState.Owner,
		teamState.PosInCurrentLeagueId,
		teamState.PosInPrevLeagueId,
		teamState.PrevLeagueId,
		teamState.InBlockIndex,
		id,
	)
	return err
}

func (b *Storage) teamHistoryAdd(id uint64, teamState TeamState) error {
	log.Infof("[DBMS] + add team history %v", teamState)
	_, err := b.db.Exec("INSERT INTO teams_history (teamId, blockNumber, currentLeagueId, owner, posInCurrentLeagueId, posInPrevLeagueId, prevLeagueId, inBlockIndex) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);",
		id,
		teamState.BlockNumber,
		teamState.CurrentLeagueId,
		teamState.Owner,
		teamState.PosInCurrentLeagueId,
		teamState.PosInPrevLeagueId,
		teamState.PrevLeagueId,
		teamState.InBlockIndex,
	)
	if err != nil {
		return err
	}

	return nil
}

func (b *Storage) TeamAdd(team Team) error {
	//  TODO: check for db is initialized
	log.Infof("[DBMS] Adding team %v", team)
	_, err := b.db.Exec("INSERT INTO teams (id, name, creationTimestamp, blockNumber, currentLeagueId, owner, posInCurrentLeagueId, posInPrevLeagueId, prevLeagueId, InBlockIndex) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);",
		team.Id,
		team.Name,
		team.CreationTimestamp,
		team.State.BlockNumber,
		team.State.CurrentLeagueId,
		team.State.Owner,
		team.State.PosInCurrentLeagueId,
		team.State.PosInPrevLeagueId,
		team.State.PrevLeagueId,
		team.State.InBlockIndex,
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
	rows, err := b.db.Query("SELECT id, name, creationTimestamp, blockNumber, currentLeagueId, owner, posInCurrentLeagueId, posInPrevLeagueId, prevLeagueId, InBlockIndex FROM teams WHERE (id = $1);", id)
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
		&team.State.InBlockIndex,
	)
	return team, nil
}
