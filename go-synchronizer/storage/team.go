package storage

import (
	log "github.com/sirupsen/logrus"
)

type TeamState struct {
	BlockNumber          string
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

func (b *Storage) TeamAdd(team Team) error {
	//  TODO: check for db is initialized
	log.Infof("(DBMS) Adding team %v %v", team.Id, team.Name)
	_, err := b.db.Exec("INSERT INTO teams (id, name, creationTimestamp) VALUES ($1, $2, $3);", team.Id, team.Name, team.CreationTimestamp)
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
	rows, err := b.db.Query("SELECT id, name FROM teams WHERE (id = $1);", id)
	if err != nil {
		return team, err
	}
	defer rows.Close()
	if !rows.Next() {
		return team, nil
	}
	rows.Scan(&team.Id, &team.Name)
	return team, nil
}
