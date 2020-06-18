package memory

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/storage"
	log "github.com/sirupsen/logrus"
)

type TeamStorageService struct {
	teams map[string]storage.Team
}

func NewTeamStorageService() *TeamStorageService {
	return &TeamStorageService{
		teams: make(map[string]storage.Team),
	}
}

func (b TeamStorageService) Team(teamId string) (*storage.Team, error) {
	team := b.teams[teamId]
	return &team, nil
}

func (b TeamStorageService) Insert(team storage.Team) error {
	b.teams[team.TeamID] = team
	return nil
}

func (b TeamStorageService) UpdateName(teamId string, name string) error {
	team, ok := b.teams[teamId]
	if !ok {
		return errors.New("unexistent team")
	}
	team.Name = name
	b.teams[teamId] = team
	return nil
}

func (b TeamStorageService) UpdateManagerName(teamId string, name string) error {
	team, ok := b.teams[teamId]
	if !ok {
		return errors.New("unexistent team")
	}
	team.ManagerName = name
	b.teams[teamId] = team
	return nil
}

func (b TeamStorageService) UpdateLeaderboardPosition(teamId string, position int) error {
	log.Warning("UpdateLeaderboardPosition not implemented")
	return nil
}
