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
	TeamUpdateZombiesFunc                     func() error
	TeamCleanZombiesFunc                      func() error
}

func (b TeamStorageService) Team(teamId string) (*storage.Team, error) {
	if b.TeamFunc == nil {
		return nil, errors.New("TeamStorageService::TeamFunc is nil")
	}
	return b.TeamFunc(teamId)
}

func (b TeamStorageService) Insert(team storage.Team) error {
	if b.InsertFunc == nil {
		return errors.New("TeamStorageService::InsertFunc is nil")
	}
	return b.InsertFunc(team)
}

func (b TeamStorageService) UpdateName(teamId string, name string) error {
	if b.UpdateNameFunc == nil {
		return errors.New("TeamStorageService::UpdateNameFunc is nil")
	}
	return b.UpdateNameFunc(teamId, name)
}

func (b TeamStorageService) UpdateManagerName(teamId string, name string) error {
	if b.UpdateManagerNameFunc == nil {
		return errors.New("TeamStorageService::UpdateManagerNameFunc is nil")
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
		return nil, errors.New("TeamStorageService::TeamsByTimezoneIdxCountryIdxLeagueIdxFunc is nil")
	}
	return b.TeamsByTimezoneIdxCountryIdxLeagueIdxFunc(timezoneIdx, countryIdx, leagueIdx)
}

func (b TeamStorageService) TeamUpdateZombies() error {
	if b.TeamUpdateZombiesFunc == nil {
		return errors.New("TeamStorageService::TeamUpdateZombies is nil")
	}
	return b.TeamUpdateZombiesFunc()
}

func (b TeamStorageService) TeamCleanZombies() error {
	if b.TeamCleanZombiesFunc == nil {
		return errors.New("TeamStorageService::TeamCleanZombies is nil")
	}
	return b.TeamCleanZombiesFunc()
}
