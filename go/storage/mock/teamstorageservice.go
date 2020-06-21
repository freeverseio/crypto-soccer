package mock

import (
	"github.com/freeverseio/crypto-soccer/go/storage"
)

type TeamStorageService struct {
	TeamFunc                      func(teamId string) (*storage.Team, error)
	InsertFunc                    func(team storage.Team) error
	UpdateNameFunc                func(teamId string, name string) error
	UpdateManagerNameFunc         func(teamId string, name string) error
	UpdateLeaderboardPositionFunc func(teamId string, position int) error
}

func (b TeamStorageService) Team(teamId string) (*storage.Team, error) {
	return b.Team(teamId)
}

func (b TeamStorageService) Insert(team storage.Team) error {
	return b.InsertFunc(team)
}

func (b TeamStorageService) UpdateName(teamId string, name string) error {
	return b.UpdateNameFunc(teamId, name)
}

func (b TeamStorageService) UpdateManagerName(teamId string, name string) error {
	return b.UpdateManagerNameFunc(teamId, name)
}

func (b TeamStorageService) UpdateLeaderboardPosition(teamId string, position int) error {
	return b.UpdateLeaderboardPositionFunc(teamId, position)
}
