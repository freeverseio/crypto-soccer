package process

import (
	"errors"
	"math/big"

	"github.com/freeverseio/crypto-soccer/go/names"
	relay "github.com/freeverseio/crypto-soccer/go/relay/storage"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/utils"
	log "github.com/sirupsen/logrus"
)

type DivisionCreationProcessor struct {
	universedb            *storage.Storage
	relaydb               *relay.Storage
	assets                *assets.Assets
	SK_SHO                uint8
	SK_SPE                uint8
	SK_PAS                uint8
	SK_DEF                uint8
	SK_END                uint8
	LEAGUES_PER_DIV       uint8
	TEAMS_PER_LEAGUE      uint8
	calendarProcessor     *Calendar
	PLAYERS_PER_TEAM_INIT uint8
}

func NewDivisionCreationProcessor(
	universedb *storage.Storage,
	relaydb *relay.Storage,
	assets *assets.Assets,
	leagues *leagues.Leagues,
) (*DivisionCreationProcessor, error) {
	SK_SHO, err := assets.SKSHO(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	SK_SPE, err := assets.SKSPE(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	SK_PAS, err := assets.SKPAS(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	SK_DEF, err := assets.SKDEF(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	SK_END, err := assets.SKEND(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	LEAGUES_PER_DIV, err := assets.LEAGUESPERDIV(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	TEAMS_PER_LEAGUE, err := assets.TEAMSPERLEAGUE(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	calendarProcessor, err := NewCalendar(leagues, universedb)
	if err != nil {
		return nil, err
	}
	PLAYERS_PER_TEAM_INIT, err := assets.PLAYERSPERTEAMINIT(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	return &DivisionCreationProcessor{
		universedb,
		relaydb,
		assets,
		SK_SHO,
		SK_SPE,
		SK_PAS,
		SK_DEF,
		SK_END,
		LEAGUES_PER_DIV,
		TEAMS_PER_LEAGUE,
		calendarProcessor,
		PLAYERS_PER_TEAM_INIT,
	}, nil
}

func (b *DivisionCreationProcessor) Process(event assets.AssetsDivisionCreation) error {
	log.Infof("Division Creation: timezoneIdx: %v, countryIdx %v, divisionIdx %v", event.Timezone, event.CountryIdxInTZ.Uint64(), event.DivisionIdxInCountry.Uint64())
	if event.CountryIdxInTZ.Uint64() == 0 {
		if err := b.universedb.TimezoneCreate(storage.Timezone{event.Timezone}); err != nil {
			return err
		}
	}
	if event.DivisionIdxInCountry.Uint64() == 0 {
		countryIdx := event.CountryIdxInTZ.Uint64()
		if countryIdx > 65535 {
			return errors.New("Cannot cast country idx to uint16: value too large")
		}
		if err := b.universedb.CountryCreate(storage.Country{event.Timezone, uint32(countryIdx)}); err != nil {
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

	leagueIdxBegin := divisionIdxInCountry.Int64() * int64(b.LEAGUES_PER_DIV)
	leagueIdxEnd := leagueIdxBegin + int64(b.LEAGUES_PER_DIV)

	for leagueIdx := leagueIdxBegin; leagueIdx < leagueIdxEnd; leagueIdx++ {
		if err := b.universedb.LeagueCreate(storage.League{timezone, uint32(countryIdx.Uint64()), uint32(leagueIdx)}); err != nil {
			return err
		}
		teamIdxBegin := leagueIdx * int64(b.TEAMS_PER_LEAGUE)
		teamIdxEnd := teamIdxBegin + int64(b.TEAMS_PER_LEAGUE)
		for teamIdxInLeague, teamIdx := uint32(0), teamIdxBegin; teamIdx < teamIdxEnd; teamIdx, teamIdxInLeague = teamIdx+1, teamIdxInLeague+1 {
			if teamId, err := b.assets.EncodeTZCountryAndVal(opts, timezone, countryIdx, big.NewInt(teamIdx)); err != nil {
				return err
			} else {
				if err := b.universedb.TeamCreate(
					storage.Team{
						teamId,
						names.GenerateTeamName(teamId),
						timezone,
						uint32(countryIdx.Uint64()),
						storage.TeamState{
							"0x0000000000000000000000000000000000000000",
							uint32(leagueIdx),
							teamIdxInLeague,
							0,
							0,
							0,
							0,
							0,
							0,
							big.NewInt(10),
							big.NewInt(0),
						},
					},
				); err != nil {
					return err
				} else if err := b.storeVirtualPlayersForTeam(opts, teamId, timezone, countryIdx, teamIdx); err != nil {
					return err
				} else if err := b.createInitialTactics(teamId); err != nil {
					return err
				}
			}
		}

		err := b.calendarProcessor.Generate(timezone, uint32(countryIdx.Uint64()), uint32(leagueIdx))
		if err != nil {
			return err
		}
		err = b.calendarProcessor.Populate(timezone, uint32(countryIdx.Uint64()), uint32(leagueIdx))
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *DivisionCreationProcessor) storeVirtualPlayersForTeam(opts *bind.CallOpts, teamId *big.Int, timezone uint8, countryIdx *big.Int, teamIdxInCountry int64) error {
	begin := teamIdxInCountry * int64(b.PLAYERS_PER_TEAM_INIT)
	end := begin + int64(b.PLAYERS_PER_TEAM_INIT)

	for i := begin; i < end; i++ {
		if playerId, err := b.assets.EncodeTZCountryAndVal(opts, timezone, countryIdx, big.NewInt(i)); err != nil {
			return err
		} else if encodedSkills, err := b.assets.GetPlayerSkillsAtBirth(opts, playerId); err != nil {
			return err
		} else if encodedState, err := b.assets.GetPlayerStateAtBirth(opts, playerId); err != nil {
			return err
		} else if defence, speed, pass, shoot, endurance, potential, dayOfBirth, err := utils.DecodeSkills(b.assets, encodedSkills); err != nil {
			return err
		} else if preferredPosition, err := b.getPlayerPreferredPosition(opts, encodedSkills); err != nil {
			return err
		} else if shirtNumber, err := b.assets.GetCurrentShirtNum(opts, encodedState); err != nil {
			return err
		} else if err := b.universedb.PlayerCreate(
			storage.Player{
				PlayerId:          playerId,
				Name:              names.GeneratePlayerName(playerId),
				PreferredPosition: preferredPosition,
				Potential:         potential.Uint64(),
				DayOfBirth:        dayOfBirth.Uint64(),
				State: storage.PlayerState{ // TODO: storage should use same skill ordering as BC
					TeamId:        teamId,
					Defence:       defence.Uint64(), // TODO: type should be uint16
					Speed:         speed.Uint64(),
					Pass:          pass.Uint64(),
					Shoot:         shoot.Uint64(),
					Endurance:     endurance.Uint64(),
					ShirtNumber:   uint8(shirtNumber.Uint64()),
					EncodedSkills: encodedSkills,
					EncodedState:  encodedState,
					Frozen:        false,
				}},
		); err != nil {
			return err
		}
	}
	return nil
}

func (p *DivisionCreationProcessor) getPlayerPreferredPosition(opts *bind.CallOpts, encodedSkills *big.Int) (string, error) {
	if forwardness, err := p.assets.GetForwardness(opts, encodedSkills); err != nil {
		return "", err
	} else if leftishness, err := p.assets.GetLeftishness(opts, encodedSkills); err != nil {
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

func (b *DivisionCreationProcessor) createInitialTactics(teamID *big.Int) error {
	tactics := b.relaydb.DefaultTactic(teamID)
	initVerse := uint64(0) // init verse
	return b.relaydb.TacticCreate(*tactics, initVerse)
}
