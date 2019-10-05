package process

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	log "github.com/sirupsen/logrus"
)

func (p *EventProcessor) StoreDivisionCreation(event leagues.LeaguesDivisionCreation) error {
	log.Infof("Division Creation: timezoneIdx: %v, countryIdx %v, divisionIdx %v", event.Timezone, event.CountryIdxInTZ.Uint64(), event.DivisionIdxInCountry.Uint64())
	if event.CountryIdxInTZ.Uint64() == 0 {
		if err := p.db.TimezoneCreate(storage.Timezone{event.Timezone}); err != nil {
			return err
		}
	}
	if event.DivisionIdxInCountry.Uint64() == 0 {
		countryIdx := event.CountryIdxInTZ.Uint64()
		if countryIdx > 65535 {
			return errors.New("Cannot cast country idx to uint16: value too large")
		}
		if err := p.db.CountryCreate(storage.Country{event.Timezone, uint32(countryIdx)}); err != nil {
			return err
		}
		if err := p.storeTeamsForNewDivision(event.Timezone, event.CountryIdxInTZ, event.DivisionIdxInCountry); err != nil {
			return err
		}
	}
	return nil
}
func (p *EventProcessor) storeTeamsForNewDivision(timezone uint8, countryIdx *big.Int, divisionIdxInCountry *big.Int) error {
	opts := &bind.CallOpts{}
	calendarProcessor, err := NewCalendar(p.leagues, p.db)
	if err != nil {
		return err
	}

	LEAGUES_PER_DIV, err := p.leagues.LEAGUESPERDIV(opts)
	if err != nil {
		return err
	}
	leagueIdxBegin := divisionIdxInCountry.Int64() * int64(LEAGUES_PER_DIV)
	leagueIdxEnd := leagueIdxBegin + int64(LEAGUES_PER_DIV)

	TEAMS_PER_LEAGUE, err := p.leagues.TEAMSPERLEAGUE(opts)
	if err != nil {
		return err
	}

	for leagueIdx := leagueIdxBegin; leagueIdx < leagueIdxEnd; leagueIdx++ {
		if err := p.db.LeagueCreate(storage.League{timezone, uint32(countryIdx.Uint64()), uint32(leagueIdx)}); err != nil {
			return err
		}
		teamIdxBegin := leagueIdx * int64(TEAMS_PER_LEAGUE)
		teamIdxEnd := teamIdxBegin + int64(TEAMS_PER_LEAGUE)
		for teamIdxInLeague, teamIdx := uint32(0), teamIdxBegin; teamIdx < teamIdxEnd; teamIdx, teamIdxInLeague = teamIdx+1, teamIdxInLeague+1 {
			if teamId, err := p.leagues.EncodeTZCountryAndVal(opts, timezone, countryIdx, big.NewInt(teamIdx)); err != nil {
				return err
			} else {
				if err := p.db.TeamCreate(
					storage.Team{
						teamId,
						timezone,
						uint32(countryIdx.Uint64()),
						storage.TeamState{"0x0", uint32(leagueIdx), teamIdxInLeague, 0, 0, 0, 0, 0, 0}},
				); err != nil {
					return err
				} else if err := p.storeVirtualPlayersForTeam(opts, teamId, timezone, countryIdx, teamIdx); err != nil {
					return err
				}
			}
		}

		err = calendarProcessor.Generate(timezone, uint32(countryIdx.Uint64()), uint32(leagueIdx))
		if err != nil {
			return err
		}
		err = calendarProcessor.Populate(timezone, uint32(countryIdx.Uint64()), uint32(leagueIdx))
		if err != nil {
			return err
		}
	}
	return err
}
