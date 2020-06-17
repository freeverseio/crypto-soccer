package memory

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/storage"
	log "github.com/sirupsen/logrus"
)

type StorageService struct {
	teams map[string]storage.Team
}

func NewStorageService() *StorageService {
	return &StorageService{
		teams: make(map[string]storage.Team),
	}
}

func (b StorageService) Team(teamId string) (*storage.Team, error) {
	team := b.teams[teamId]
	return &team, nil
}

func (b StorageService) Insert(team storage.Team) error {
	b.teams[team.TeamID] = team
	return nil
}

func (b StorageService) UpdateName(teamId string, name string) error {
	team, ok := b.teams[teamId]
	if !ok {
		return errors.New("unexistent team")
	}
	team.Name = name
	b.teams[teamId] = team
	return nil
}

func (b StorageService) UpdateManagerName(teamId string, name string) error {
	team, ok := b.teams[teamId]
	if !ok {
		return errors.New("unexistent team")
	}
	team.ManagerName = name
	b.teams[teamId] = team
	return nil
}

func (b StorageService) Dump(fileName string) error {
	log.Warning("Dump not implemented")
	return nil
}
