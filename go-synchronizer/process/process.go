package process

import (
	"context"
	"errors"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/updates"

	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	//"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/utils"
	log "github.com/sirupsen/logrus"
)

type EventProcessor struct {
	usesGanache bool
	client      *ethclient.Client
	db          *storage.Storage
	engine      *engine.Engine
	leagues     *leagues.Leagues
	updates     *updates.Updates
}

// *****************************************************************************
// public
// *****************************************************************************

// NewEventProcessor creates a new struct for scanning and storing crypto soccer events
func NewEventProcessor(client *ethclient.Client, db *storage.Storage, engine *engine.Engine, leagues *leagues.Leagues, updates *updates.Updates) *EventProcessor {
	return &EventProcessor{false, client, db, engine, leagues, updates}
}

// NewGanacheEventProcessor creates a new struct for scanning and storing crypto soccer events from a ganache client
func NewGanacheEventProcessor(client *ethclient.Client, db *storage.Storage, engine *engine.Engine, leagues *leagues.Leagues, updates *updates.Updates) *EventProcessor {
	return &EventProcessor{true, client, db, engine, leagues, updates}
}

// Process processes all scanned events and stores them into the database db
func (p *EventProcessor) Process(delta uint64) (uint64, error) {
	opts, err := p.nextRange(delta)
	if err != nil {
		return 0, err
	}

	if opts == nil {
		log.Info("No new blocks to scan.")
		return 0, nil
	}

	log.WithFields(log.Fields{
		"start": opts.Start,
		"end":   *opts.End,
	}).Info("Syncing ...")

	divisionCreationIter, err := p.leagues.FilterDivisionCreation(opts)
	if err != nil {
		return 0, err
	}

	teamTransferIter, err := p.leagues.FilterTeamTransfer(opts)
	if err != nil {
		return 0, err
	}

	playerTransferIter, err := p.leagues.FilterPlayerTransfer(opts)
	if err != nil {
		return 0, err
	}

	actionSubmissionIter, err := p.updates.FilterActionsSubmission(opts)
	if err != nil {
		return 0, err
	}

	scanner := NewEventScanner(p.leagues, p.updates)
	if err := scanner.ScanActionsSubmission(actionSubmissionIter); err != nil {
		return 0, err
	}
	if err := scanner.ScanDivisionCreation(divisionCreationIter); err != nil {
		return 0, err
	}
	if err := scanner.ScanTeamTransfer(teamTransferIter); err != nil {
		return 0, err
	}
	if err := scanner.ScanPlayerTransfer(playerTransferIter); err != nil {
		return 0, err
	}
	if err := scanner.Process(opts); err != nil {
		return 0, err
	} else {
		log.Debug("scanner got: ", len(scanner.Events), " Abstract Events")
	}

	for _, v := range scanner.Events {
		if err := p.dispatch(v); err != nil {
			return 0, err
		}
	}

	deltaBlock := *opts.End - opts.Start

	// store the last block that was scanned
	return deltaBlock, p.db.SetBlockNumber(*opts.End)
}

// *****************************************************************************
// private
// *****************************************************************************
func (p *EventProcessor) dispatch(e *AbstractEvent) error {
	switch v := e.Value.(type) {
	case leagues.LeaguesDivisionCreation:
		log.Info("[processor] Dispatching LeaguesDivisionCreation event")
		return p.storeDivisionCreation(v)
	case leagues.LeaguesTeamTransfer:
		log.Info("[processor] dispatching LeaguesTeamTransfer event")
		teamID := v.TeamId
		newOwner := v.To.String()
		team, err := p.db.GetTeam(teamID)
		if err != nil {
			return err
		}
		// team.State.BlockNumber = blockNumber
		team.State.Owner = newOwner
		return p.db.TeamUpdate(teamID, team.State)
	case leagues.LeaguesPlayerTransfer:
		log.Info("[processor] dispatching LeaguesPlayerTransfer event")
		playerID := v.PlayerId
		toTeamID := v.TeamIdTarget
		player, err := p.db.GetPlayer(playerID)
		if err != nil {
			return err
		}
		playerState, err := p.leagues.GetPlayerState(&bind.CallOpts{}, playerID)
		if err != nil {
			return err
		}
		shirtNumber, err := p.leagues.GetCurrentShirtNum(&bind.CallOpts{}, playerState)
		if err != nil {
			return err
		}
		player.State.TeamId = toTeamID
		player.State.ShirtNumber = uint8(shirtNumber.Uint64())
		return p.db.PlayerUpdate(playerID, player.State)
	case updates.UpdatesActionsSubmission:
		log.Info("[processor] Dispatching UpdatesActionsSubmission event")
		leagueProcessor, err := NewLeagueProcessor(p.engine, p.leagues, p.db)
		if err != nil {
			return err
		}
		return leagueProcessor.Process(v)
	}
	return fmt.Errorf("[processor] Error dispatching unknown event type: %s", e.Name)
}
func (p *EventProcessor) nextRange(delta uint64) (*bind.FilterOpts, error) {
	start, err := p.dbLastBlockNumber()
	if err != nil {
		return nil, err
	}
	if start != 0 {
		// unless this is the very first execution,
		// the block number that is stored in the db
		// was already scanned. We are interested in
		// the next block
		if start < math.MaxUint64 {
			start += 1
		} else {
			return nil, errors.New("Block range overflow")
		}
	}
	end := p.clientLastBlockNumber()
	if delta != 0 {
		end = uint64(math.Min(float64(start+delta), float64(end)))
	}
	if start > end {
		return nil, nil
	}
	return &bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	}, nil
}

func (p *EventProcessor) clientLastBlockNumber() uint64 {
	if p.client == nil {
		log.Warn("Client is nil. Returning 0 as last block.")
		return 0
	}
	header, err := p.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Warn("Could not get blockchain last block")
		return 0
	}
	return header.Number.Uint64()
}
func (p *EventProcessor) dbLastBlockNumber() (uint64, error) {
	storedLastBlockNumber, err := p.db.GetBlockNumber()
	if err != nil {
		return 0, err
	}
	return storedLastBlockNumber, err
}
func (p *EventProcessor) getTimeOfEvent(eventRaw types.Log) (uint64, uint64, error) {
	if p.usesGanache {
		return eventRaw.BlockNumber, eventRaw.BlockNumber, nil
	}
	block, err := p.client.BlockByHash(context.Background(), eventRaw.BlockHash)
	if err != nil {
		return 0, 0, err
	}
	return block.Time(), eventRaw.BlockNumber, nil
}

func (p *EventProcessor) storeDivisionCreation(event leagues.LeaguesDivisionCreation) error {
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
				if teamOwner, err := p.leagues.GetOwnerTeam(opts, teamId); err != nil {
					return err
				} else if err := p.db.TeamCreate(
					storage.Team{
						teamId,
						timezone,
						uint32(countryIdx.Uint64()),
						storage.TeamState{teamOwner.Hex(), uint32(leagueIdx), teamIdxInLeague, 0, 0, 0, 0, 0, 0}},
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
func (p *EventProcessor) storeVirtualPlayersForTeam(opts *bind.CallOpts, teamId *big.Int, timezone uint8, countryIdx *big.Int, teamIdxInCountry int64) error {
	PLAYERS_PER_TEAM_INIT, err := p.leagues.PLAYERSPERTEAMINIT(opts)
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

	SK_SHO, err = p.leagues.SKSHO(opts)
	if err != nil {
		return err
	}
	SK_SPE, err = p.leagues.SKSPE(opts)
	if err != nil {
		return err
	}
	SK_PAS, err = p.leagues.SKPAS(opts)
	if err != nil {
		return err
	}
	SK_DEF, err = p.leagues.SKDEF(opts)
	if err != nil {
		return err
	}
	SK_END, err = p.leagues.SKEND(opts)
	if err != nil {
		return err
	}

	for i := begin; i < end; i++ {
		if playerId, err := p.leagues.EncodeTZCountryAndVal(opts, timezone, countryIdx, big.NewInt(i)); err != nil {
			return err
		} else if skills, err := p.getPlayerSkillsAtBirth(opts, playerId); err != nil {
			return err
		} else if preferredPosition, err := p.getPlayerPreferredPosition(opts, playerId); err != nil {
			return err
		} else if shirtNumber, err := p.getShirtNumber(opts, playerId); err != nil {
			return err
		} else if err := p.db.PlayerCreate(
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

func (p *EventProcessor) getPlayerSkillsAtBirth(opts *bind.CallOpts, playerId *big.Int) ([5]uint16, error) {
	if skills, err := p.leagues.GetPlayerSkillsAtBirth(opts, playerId); err != nil {
		return [5]uint16{0, 0, 0, 0, 0}, err
	} else {
		return p.leagues.GetSkillsVec(opts, skills)
	}
}

func (p *EventProcessor) getShirtNumber(opts *bind.CallOpts, playerId *big.Int) (uint8, error) {
	if playerState, err := p.leagues.GetPlayerState(opts, playerId); err != nil {
		return 0, err
	} else if shirtNumber, err := p.leagues.GetCurrentShirtNum(opts, playerState); err != nil {
		return 0, err
	} else {
		return uint8(shirtNumber.Uint64()), nil
	}
}

func (p *EventProcessor) getPlayerPreferredPosition(opts *bind.CallOpts, playerId *big.Int) (string, error) {
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
