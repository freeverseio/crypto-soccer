package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/storage"
)

type StorageService struct {
	tx *sql.Tx
}

func NewStorageService(tx *sql.Tx) *StorageService {
	return &StorageService{
		tx: tx,
	}
}

func (b StorageService) Team(teamId string) (*storage.Team, error) {
	team, err := storage.TeamByTeamId(b.tx, teamId)
	return &team, err
}

func (b StorageService) Insert(team storage.Team) error {
	return team.Insert(b.tx)
}

func (b StorageService) UpdateName(teamId string, name string) error {
	team, err := storage.TeamByTeamId(b.tx, teamId)
	if err != nil {
		return err
	}
	team.Name = name
	return team.Update(b.tx)
}

func (b StorageService) UpdateManagerName(teamId string, name string) error {
	team, err := storage.TeamByTeamId(b.tx, teamId)
	if err != nil {
		return err
	}
	team.ManagerName = name
	return team.Update(b.tx)
}
