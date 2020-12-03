package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/storage"
)

type TeamStorageService struct {
	tx *sql.Tx
}

func NewTeamStorageService(tx *sql.Tx) *TeamStorageService {
	return &TeamStorageService{
		tx: tx,
	}
}

func (b TeamStorageService) Team(teamId string) (*storage.Team, error) {
	team, err := storage.TeamByTeamId(b.tx, teamId)
	return &team, err
}

func (b TeamStorageService) Insert(team storage.Team) error {
	return team.Insert(b.tx)
}

func (b TeamStorageService) UpdateName(teamId string, name string) error {
	team, err := storage.TeamByTeamId(b.tx, teamId)
	if err != nil {
		return err
	}
	team.Name = name
	return team.Update(b.tx)
}

func (b TeamStorageService) UpdateManagerName(teamId string, name string) error {
	team, err := storage.TeamByTeamId(b.tx, teamId)
	if err != nil {
		return err
	}
	team.ManagerName = name
	return team.Update(b.tx)
}

func (b TeamStorageService) UpdateLeaderboardPosition(teamId string, position int) error {
	team, err := storage.TeamByTeamId(b.tx, teamId)
	if err != nil {
		return err
	}
	team.LeaderboardPosition = position
	return team.Update(b.tx)
}

func (b TeamStorageService) TeamsByTimezoneIdxCountryIdxLeagueIdx(timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) ([]storage.Team, error) {
	return storage.TeamsByTimezoneIdxCountryIdxLeagueIdx(b.tx, timezoneIdx, countryIdx, leagueIdx)
}

func (b TeamStorageService) TeamUpdateZombies(timezoneIdx uint8, countryIdx uint32) error {

	return storage.TeamUpdateZombies(b.tx, timezoneIdx, countryIdx)
}

func (b TeamStorageService) TeamCleanZombies(timezoneIdx uint8, countryIdx uint32) error {

	return storage.TeamCleanZombies(b.tx, timezoneIdx, countryIdx)
}

func (b TeamStorageService) TeamPromoTimeout(teamId string) (uint32, error) {
	team, err := storage.TeamByTeamId(b.tx, teamId)
	if err != nil {
		return 0, err
	}
	return team.PromoTimeout, nil
}

func (b TeamStorageService) TeamSetPromoTimeout(teamId string, promoTimeout uint32) error {
	_, err := b.tx.Exec("UPDATE teams SET promo_timeout=$1 WHERE team_id=$2", promoTimeout, teamId)
	return err
}
