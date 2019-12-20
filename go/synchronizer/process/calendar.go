package process

import (
	"database/sql"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

type Calendar struct {
	contracts   *contracts.Contracts
	MatchDays   uint8
	MatchPerDay uint8
}

func NewCalendar(contracts *contracts.Contracts) (*Calendar, error) {
	matchDays, err := contracts.Leagues.MATCHDAYS(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	matchPerDay, err := contracts.Leagues.MATCHESPERDAY(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	return &Calendar{contracts, matchDays, matchPerDay}, nil
}

func (b *Calendar) Generate(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) error {
	league, err := storage.GetLeague(tx, leagueIdx)
	if err != nil {
		return err
	}
	if league == nil {
		return errors.New("Unexistent league")
	}

	for matchDay := uint8(0); matchDay < b.MatchDays; matchDay++ {
		for match := uint8(0); match < b.MatchPerDay; match++ {
			match := storage.Match{
				TimezoneIdx: timezoneIdx,
				CountryIdx:  countryIdx,
				LeagueIdx:   leagueIdx,
				MatchDayIdx: matchDay,
				MatchIdx:    match,
			}
			err = match.MatchCreate(tx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Calendar) Populate(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) error {
	league, err := storage.GetLeague(tx, leagueIdx)
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
			homeTeamID, err := storage.TeamIdByTimezoneIdxCountryIdxLeagueIdx(tx, timezoneIdx, countryIdx, leagueIdx, uint32(teams.HomeIdx))
			if err != nil {
				return err
			}
			visitorTeamID, err := storage.TeamIdByTimezoneIdxCountryIdxLeagueIdx(tx, timezoneIdx, countryIdx, leagueIdx, uint32(teams.VisitorIdx))
			if err != nil {
				return err
			}
			err = storage.MatchSetTeams(tx, timezoneIdx, countryIdx, leagueIdx, matchDay, match, homeTeamID, visitorTeamID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Calendar) Reset(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) error {
	league, err := storage.GetLeague(tx, leagueIdx)
	if err != nil {
		return err
	}
	if league == nil {
		return errors.New("Unexistent league")
	}

	for matchDay := uint8(0); matchDay < b.MatchDays; matchDay++ {
		for match := uint8(0); match < b.MatchPerDay; match++ {
			err = storage.MatchReset(tx, timezoneIdx, countryIdx, leagueIdx, matchDay, match)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
