package mock

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/storage"
)

type TeamStorageService struct {
	TeamFunc                                  func(teamId string) (*storage.Team, error)
	InsertFunc                                func(team storage.Team) error
	UpdateNameFunc                            func(teamId string, name string) error
	UpdateManagerNameFunc                     func(teamId string, name string) error
	UpdateLeaderboardPositionFunc             func(teamId string, position int) error
	TeamsByTimezoneIdxCountryIdxLeagueIdxFunc func(timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) ([]storage.Team, error)
}

func (b TeamStorageService) Team(teamId string) (*storage.Team, error) {
	if b.TeamFunc == nil {
		return nil, errors.New("TeamStorageService::Team is nil")
	}
	return b.TeamFunc(teamId)
}

func (b TeamStorageService) Insert(team storage.Team) error {
	if b.InsertFunc == nil {
		return errors.New("TeamStorageService::Insert is nil")
	}
	return b.InsertFunc(team)
}

func (b TeamStorageService) UpdateName(teamId string, name string) error {
	if b.UpdateNameFunc == nil {
		return errors.New("TeamStorageService::UpdateName is nil")
	}
	return b.UpdateNameFunc(teamId, name)
}

func (b TeamStorageService) UpdateManagerName(teamId string, name string) error {
	if b.UpdateManagerNameFunc == nil {
		return errors.New("TeamStorageService::UpdateManagerName is nil")
	}
	return b.UpdateManagerNameFunc(teamId, name)
}

func (b TeamStorageService) UpdateLeaderboardPosition(teamId string, position int) error {
	if b.UpdateLeaderboardPositionFunc == nil {
		return errors.New("TeamStorageService::UpdateLeaderboardPositionFunc is nil")
	}
	return b.UpdateLeaderboardPositionFunc(teamId, position)
}

func (b TeamStorageService) TeamsByTimezoneIdxCountryIdxLeagueIdx(timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) ([]storage.Team, error) {
	if b.TeamsByTimezoneIdxCountryIdxLeagueIdxFunc == nil {
		return nil, errors.New("TeamStorageService::TeamsByTimezoneIdxCountryIdxLeagueIdx is nil")
	}
	return b.TeamsByTimezoneIdxCountryIdxLeagueIdxFunc(timezoneIdx, countryIdx, leagueIdx)
}
