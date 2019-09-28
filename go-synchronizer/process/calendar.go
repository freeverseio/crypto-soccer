package process

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
)

type Calendar struct {
	leagues *leagues.Leagues
	storage *storage.Storage
}

func NewCalendar(leagues *leagues.Leagues, storage *storage.Storage) *Calendar {
	return &Calendar{leagues, storage}
}

func (b *Calendar) Generate(timezoneIdx uint8, countryIdx uint16, leagueIdx uint32) error {
	league, err := b.storage.GetLeague(leagueIdx)
	if err != nil {
		return err
	}
	if league == nil {
		return errors.New("Unexistent league")
	}
	return nil
}
