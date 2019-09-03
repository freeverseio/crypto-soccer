package process

import (
	"context"
	//"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/states"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	log "github.com/sirupsen/logrus"
)

type EventProcessor struct {
	usesGanache bool
	client      *ethclient.Client
	db          *storage.Storage
	assets      *assets.Assets
	states      *states.States
	leagues     *leagues.Leagues
}

// *****************************************************************************
// public
// *****************************************************************************

// NewEventProcessor creates a new struct for scanning and storing crypto soccer events
func NewEventProcessor(client *ethclient.Client, db *storage.Storage, assets *assets.Assets, states *states.States, leagues *leagues.Leagues) *EventProcessor {
	return &EventProcessor{false, client, db, assets, states, leagues}
}

// NewGanacheEventProcessor creates a new struct for scanning and storing crypto soccer events from a ganache client
func NewGanacheEventProcessor(client *ethclient.Client, db *storage.Storage, assets *assets.Assets, states *states.States, leagues *leagues.Leagues) *EventProcessor {
	return &EventProcessor{true, client, db, assets, states, leagues}
}

// Process processes all scanned events and stores them into the database db
func (p *EventProcessor) Process() error {
	opts := p.nextRange()

	if opts == nil {
		log.Info("No new blocks to scan.")
		return nil
	}

	log.WithFields(log.Fields{
		"start": opts.Start,
		"end":   *opts.End,
	}).Info("Syncing ...")

	if events, err := p.scanTeamCreated(opts); err != nil {
		return err
	} else {
		err = p.storeTeamCreated(events)
		if err != nil {
			return err
		}
	}

	if events, err := p.scanTeamTransfer(opts); err != nil {
		return err
	} else {
		for _, event := range events { // TODO: next part to be recoded
			_, blockNumber, err := p.getTimeOfEvent(event.Raw)
			if err != nil {
				return err
			}
			teamId := event.TeamId.Uint64()
			newOwner := event.To.String()
			team, err := p.db.GetTeam(teamId)
			if err != nil {
				return err
			}
			team.State.BlockNumber = blockNumber
			team.State.Owner = newOwner
			err = p.db.TeamStateUpdate(teamId, team.State)
			if err != nil {
				return err
			}
		}
	}

	if p.leagues != nil {
		if events, err := p.scanLeagueCreated(opts); err != nil {
			return err
		} else {
			for _, event := range events { // TODO: next part to be recoded
				p.db.LeagueAdd(storage.League{
					Id: event.LeagueId.Uint64(),
				})
				// log.Info(
				// 	"Found league ", event.LeagueId.Int64(),
				// 	"\n\tdays: ", p.getLeagueDaysCount(event.LeagueId),
				// 	"\n\tfinished: ", p.hasLeagueFinished(event.LeagueId),
				// 	"\n\tupdated: ", p.isLeagueUpdated(event.LeagueId),
				// )
			}
		}
	} else {
		log.Warn("Contract leagues not set. Skipping scanning for leagues")
	}

	// store the last block that was scanned
	p.db.SetBlockNumber(*opts.End)
	return nil
}

// *****************************************************************************
// private
// *****************************************************************************
func (p *EventProcessor) nextRange() *bind.FilterOpts {
	start := p.dbLastBlockNumber()
	if start != 0 {
		// unless this is the very first execution,
		// the block number that is stored in the db
		// was already scanned. We are interested in
		// the next block
		if start < math.MaxUint64 {
			start += 1
		} else {
			log.Error("Block range overflow")
			return nil
		}
	}
	end := p.clientLastBlockNumber()
	if start > end {
		return nil
	}
	return &bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	}
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
func (p *EventProcessor) dbLastBlockNumber() uint64 {
	storedLastBlockNumber, err := p.db.GetBlockNumber()
	if err != nil {
		log.Warn("Could not get database last block")
		return 0
	}
	return storedLastBlockNumber
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
func (p *EventProcessor) storeTeamCreated(events []assets.AssetsTeamCreated) error {
	for _, event := range events {
		if name, err := p.assets.GetTeamName(nil, event.Id); err != nil {
			return err
		} else if owner, err := p.assets.GetTeamOwner(nil, event.Id); err != nil {
			return err
		} else if blockTime, blockNumber, err := p.getTimeOfEvent(event.Raw); err != nil {
			return err
		} else if err := p.db.TeamAdd(storage.Team{
			Id:                event.Id.Uint64(),
			Name:              name,
			CreationTimestamp: blockTime,
			CountryId:         1, // TODO: get it from blockchain
			State: storage.TeamState{
				BlockNumber:          blockNumber,
				Owner:                owner.Hex(),
				CurrentLeagueId:      1, // TODO: uint64
				PosInCurrentLeagueId: 0, // TODO: uint64
				PrevLeagueId:         0, // TODO: uint64
				PosInPrevLeagueId:    0, // TODO: uint64
			},
		}); err != nil {
			return err
		}
		if err := p.storeVirtualPlayers(event.Id); err != nil {
			return err
		}
	}
	return nil
}
func (p *EventProcessor) scanTeamCreated(opts *bind.FilterOpts) ([]assets.AssetsTeamCreated, error) {
	if opts == nil {
		opts = &bind.FilterOpts{Start: 0}
	}
	iter, err := p.assets.FilterTeamCreated(opts)
	if err != nil {
		return nil, err
	}

	events := []assets.AssetsTeamCreated{}

	for iter.Next() {
		events = append(events, *(iter.Event))
	}
	return events, nil
}
func (p *EventProcessor) scanTeamTransfer(opts *bind.FilterOpts) ([]assets.AssetsTeamTransfer, error) {
	if opts == nil {
		opts = &bind.FilterOpts{Start: 0}
	}
	iter, err := p.assets.FilterTeamTransfer(opts)
	if err != nil {
		return nil, err
	}

	events := []assets.AssetsTeamTransfer{}

	for iter.Next() {
		events = append(events, *(iter.Event))
	}
	return events, nil
}
func (p *EventProcessor) storeVirtualPlayers(teamId *big.Int) error {
	// TODO: move to a single run place ...  constructor
	nPlayersAtCreation, err := p.assets.PLAYERSPERTEAMINIT(&bind.CallOpts{})
	if err != nil {
		return err
	}

	for i := 0; i < int(nPlayersAtCreation); i++ {
		if id, err := p.assets.GenerateVirtualPlayerId(&bind.CallOpts{}, teamId, uint8(i)); err != nil {
			return err
		} else if state, err := p.assets.GenerateVirtualPlayerState(&bind.CallOpts{}, id); err != nil {
			return err
		} else {
			if skills, err := p.states.GetSkillsVec(&bind.CallOpts{}, state); err != nil {
				return err
			} else {
				player := storage.Player{
					Id:                     id.Uint64(),
					MonthOfBirthInUnixTime: "0", // TODO
					State: storage.PlayerState{
						TeamId:    teamId.Uint64(),
						State:     state.String(),
						Defence:   uint64(skills[0]),
						Speed:     uint64(skills[1]),
						Pass:      uint64(skills[2]),
						Shoot:     uint64(skills[3]),
						Endurance: uint64(skills[4]),
					},
				}
				p.db.PlayerAdd(player)
			}
		}
	}
	return nil
}
func (p *EventProcessor) scanLeagueCreated(opts *bind.FilterOpts) ([]leagues.LeaguesLeagueCreated, error) {
	if opts == nil {
		opts = &bind.FilterOpts{Start: 0}
	}
	iter, err := p.leagues.FilterLeagueCreated(opts)
	if err != nil {
		return nil, err
	}

	events := []leagues.LeaguesLeagueCreated{}

	for iter.Next() {
		events = append(events, *(iter.Event))
	}
	return events, nil
}
func (p *EventProcessor) isLeagueUpdated(leagueId *big.Int) bool {
	result, err := p.leagues.IsUpdated(nil, leagueId)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return result
}
func (p *EventProcessor) hasLeagueFinished(leagueId *big.Int) bool {
	result, err := p.leagues.HasFinished(nil, leagueId)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return result
}
func (p *EventProcessor) getLeagueDaysCount(leagueId *big.Int) int64 {
	result, err := p.leagues.CountLeagueDays(nil, leagueId)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	return result.Int64()
}
