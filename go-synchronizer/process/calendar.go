package process

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
)

type Calendar struct {
	leagues     *leagues.Leagues
	storage     *storage.Storage
	MatchDays   uint8
	MatchPerDay uint8
}

func NewCalendar(leagues *leagues.Leagues, storage *storage.Storage) (*Calendar, error) {
	matchDays, err := leagues.MATCHDAYS(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	matchPerDay, err := leagues.MATCHESPERDAY(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	return &Calendar{leagues, storage, matchDays, matchPerDay}, nil
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
