package process

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/utils"
	log "github.com/sirupsen/logrus"
)

type DivisionCreationProcessor struct {
	db      *storage.Storage
	leagues *leagues.Leagues
}

func NewDivisionCreationProcessor(db *storage.Storage, leagues *leagues.Leagues) *DivisionCreationProcessor {
	return &DivisionCreationProcessor{db, leagues}
}

func (b *DivisionCreationProcessor) StoreDivisionCreation(event leagues.LeaguesDivisionCreation) error {
	log.Infof("Division Creation: timezoneIdx: %v, countryIdx %v, divisionIdx %v", event.Timezone, event.CountryIdxInTZ.Uint64(), event.DivisionIdxInCountry.Uint64())
	if event.CountryIdxInTZ.Uint64() == 0 {
		if err := b.db.TimezoneCreate(storage.Timezone{event.Timezone}); err != nil {
			return err
		}
	}
	if event.DivisionIdxInCountry.Uint64() == 0 {
		countryIdx := event.CountryIdxInTZ.Uint64()
		if countryIdx > 65535 {
			return errors.New("Cannot cast country idx to uint16: value too large")
		}
		if err := b.db.CountryCreate(storage.Country{event.Timezone, uint32(countryIdx)}); err != nil {
			return err
		}
		if err := b.storeTeamsForNewDivision(event.Timezone, event.CountryIdxInTZ, event.DivisionIdxInCountry); err != nil {
			return err
		}
	}
	return nil
}
func (b *DivisionCreationProcessor) storeTeamsForNewDivision(timezone uint8, countryIdx *big.Int, divisionIdxInCountry *big.Int) error {
	opts := &bind.CallOpts{}
	calendarProcessor, err := NewCalendar(b.leagues, b.db)
	if err != nil {
		return err
	}

	LEAGUES_PER_DIV, err := b.leagues.LEAGUESPERDIV(opts)
	if err != nil {
		return err
	}
	leagueIdxBegin := divisionIdxInCountry.Int64() * int64(LEAGUES_PER_DIV)
	leagueIdxEnd := leagueIdxBegin + int64(LEAGUES_PER_DIV)

	TEAMS_PER_LEAGUE, err := b.leagues.TEAMSPERLEAGUE(opts)
	if err != nil {
		return err
	}

	for leagueIdx := leagueIdxBegin; leagueIdx < leagueIdxEnd; leagueIdx++ {
		if err := b.db.LeagueCreate(storage.League{timezone, uint32(countryIdx.Uint64()), uint32(leagueIdx)}); err != nil {
			return err
		}
		teamIdxBegin := leagueIdx * int64(TEAMS_PER_LEAGUE)
		teamIdxEnd := teamIdxBegin + int64(TEAMS_PER_LEAGUE)
		for teamIdxInLeague, teamIdx := uint32(0), teamIdxBegin; teamIdx < teamIdxEnd; teamIdx, teamIdxInLeague = teamIdx+1, teamIdxInLeague+1 {
			if teamId, err := b.leagues.EncodeTZCountryAndVal(opts, timezone, countryIdx, big.NewInt(teamIdx)); err != nil {
				return err
			} else {
				if err := b.db.TeamCreate(
					storage.Team{
						teamId,
						timezone,
						uint32(countryIdx.Uint64()),
						storage.TeamState{"0x0", uint32(leagueIdx), teamIdxInLeague, 0, 0, 0, 0, 0, 0}},
				); err != nil {
					return err
				} else if err := b.storeVirtualPlayersForTeam(opts, teamId, timezone, countryIdx, teamIdx); err != nil {
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

func (b *DivisionCreationProcessor) storeVirtualPlayersForTeam(opts *bind.CallOpts, teamId *big.Int, timezone uint8, countryIdx *big.Int, teamIdxInCountry int64) error {
	PLAYERS_PER_TEAM_INIT, err := b.leagues.PLAYERSPERTEAMINIT(opts)
	if err != nil {
		return err
	}
	begin := teamIdxInCountry * int64(PLAYERS_PER_TEAM_INIT)
	end := begin + int64(PLAYERS_PER_TEAM_INIT)

	SK_SHO := uint8(0)
	SK_SPE := uint8(0)
	SK_PAS := uint8(0)
	SK_DEF := uint8(0)
	SK_END := uint8(0)

	SK_SHO, err = b.leagues.SKSHO(opts)
	if err != nil {
		return err
	}
	SK_SPE, err = b.leagues.SKSPE(opts)
	if err != nil {
		return err
	}
	SK_PAS, err = b.leagues.SKPAS(opts)
	if err != nil {
		return err
	}
	SK_DEF, err = b.leagues.SKDEF(opts)
	if err != nil {
		return err
	}
	SK_END, err = b.leagues.SKEND(opts)
	if err != nil {
		return err
	}

	for i := begin; i < end; i++ {
		if playerId, err := b.leagues.EncodeTZCountryAndVal(opts, timezone, countryIdx, big.NewInt(i)); err != nil {
			return err
		} else if skills, err := b.getPlayerSkillsAtBirth(opts, playerId); err != nil {
			return err
		} else if preferredPosition, err := b.getPlayerPreferredPosition(opts, playerId); err != nil {
			return err
		} else if shirtNumber, err := b.getShirtNumber(opts, playerId); err != nil {
			return err
		} else if err := b.db.PlayerCreate(
			storage.Player{
				playerId,
				preferredPosition,
				storage.PlayerState{ // TODO: storage should use same skill ordering as BC
					TeamId:      teamId,
					Defence:     uint64(skills[SK_DEF]), // TODO: type should be uint16
					Speed:       uint64(skills[SK_SPE]),
					Pass:        uint64(skills[SK_PAS]),
					Shoot:       uint64(skills[SK_SHO]),
					Endurance:   uint64(skills[SK_END]),
					ShirtNumber: shirtNumber,
				}},
		); err != nil {
			return err
		}
	}
	return err
}

func (p *DivisionCreationProcessor) getPlayerSkillsAtBirth(opts *bind.CallOpts, playerId *big.Int) ([5]uint16, error) {
	if skills, err := p.leagues.GetPlayerSkillsAtBirth(opts, playerId); err != nil {
		return [5]uint16{0, 0, 0, 0, 0}, err
	} else {
		return p.leagues.GetSkillsVec(opts, skills)
	}
}

func (p *DivisionCreationProcessor) getShirtNumber(opts *bind.CallOpts, playerId *big.Int) (uint8, error) {
	if playerState, err := p.leagues.GetPlayerState(opts, playerId); err != nil {
		return 0, err
	} else if shirtNumber, err := p.leagues.GetCurrentShirtNum(opts, playerState); err != nil {
		return 0, err
	} else {
		return uint8(shirtNumber.Uint64()), nil
	}
}

func (p *DivisionCreationProcessor) getPlayerPreferredPosition(opts *bind.CallOpts, playerId *big.Int) (string, error) {
	if encodedSkills, err := p.leagues.GetPlayerSkillsAtBirth(opts, playerId); err != nil {
		return "", err
	} else if forwardness, err := p.leagues.GetForwardness(opts, encodedSkills); err != nil {
		return "", err
	} else if leftishness, err := p.leagues.GetLeftishness(opts, encodedSkills); err != nil {
		return "", err
	} else {
		if forwardness.Uint64() > 255 {
			return "", errors.New("Cannot cast forwardness to uint8: value too large")
		} else if leftishness.Uint64() > 255 {
			return "", errors.New("Cannot cast leftishness to uint8: value too large")
		}
		return utils.PreferredPosition(uint8(forwardness.Uint64()), uint8(leftishness.Uint64()))
	}
}
