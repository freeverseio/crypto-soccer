package process

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

type Calendar struct {
	contracts   *contracts.Contracts
	storage     *storage.Storage
	MatchDays   uint8
	MatchPerDay uint8
}

func NewCalendar(contracts *contracts.Contracts, storage *storage.Storage) (*Calendar, error) {
	matchDays, err := contracts.Leagues.MATCHDAYS(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	matchPerDay, err := contracts.Leagues.MATCHESPERDAY(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	return &Calendar{contracts, storage, matchDays, matchPerDay}, nil
}

func (b *Calendar) Generate(timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) error {
	league, err := b.storage.GetLeague(leagueIdx)
	if err != nil {
		return err
	}
	if league == nil {
		return errors.New("Unexistent league")
	}

	for matchDay := uint8(0); matchDay < b.MatchDays; matchDay++ {
		for match := uint8(0); match < b.MatchPerDay; match++ {
			err = b.storage.MatchCreate(storage.Match{
				TimezoneIdx: timezoneIdx,
				CountryIdx:  countryIdx,
				LeagueIdx:   leagueIdx,
				MatchDayIdx: matchDay,
				MatchIdx:    match,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Calendar) Populate(timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) error {
	league, err := b.storage.GetLeague(leagueIdx)
	if err != nil {
		return err
	}
	if league == nil {
		return errors.New("Unexistent league")
	}

	for matchDay := uint8(0); matchDay < b.MatchDays; matchDay++ {
		for match := uint8(0); match < b.MatchPerDay; match++ {
			teams, err := b.contracts.Leagues.GetTeamsInLeagueMatch(&bind.CallOpts{}, matchDay, match)
			if err != nil {
				return nil
			}
			homeTeamID, err := b.storage.GetTeamID(timezoneIdx, countryIdx, leagueIdx, uint32(teams.HomeIdx))
			if err != nil {
				return err
			}
			visitorTeamID, err := b.storage.GetTeamID(timezoneIdx, countryIdx, leagueIdx, uint32(teams.VisitorIdx))
			if err != nil {
				return err
			}
			err = b.storage.MatchSetTeams(timezoneIdx, countryIdx, leagueIdx, matchDay, match, homeTeamID, visitorTeamID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Calendar) Reset(timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) error {
	league, err := b.storage.GetLeague(leagueIdx)
	if err != nil {
		return err
	}
	if league == nil {
		return errors.New("Unexistent league")
	}

	for matchDay := uint8(0); matchDay < b.MatchDays; matchDay++ {
		for match := uint8(0); match < b.MatchPerDay; match++ {
			err = b.storage.MatchReset(timezoneIdx, countryIdx, leagueIdx, matchDay, match)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
