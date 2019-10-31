package process

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/freeverseio/crypto-soccer/go/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go/contracts/market"
	"github.com/freeverseio/crypto-soccer/go/contracts/updates"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	log "github.com/sirupsen/logrus"
)

type EventProcessor struct {
	client                    *ethclient.Client
	db                        *storage.Storage
	engine                    *engine.Engine
	assets                    *assets.Assets
	leagues                   *leagues.Leagues
	updates                   *updates.Updates
	market                    *market.Market
	divisionCreationProcessor *DivisionCreationProcessor
	leagueProcessor           *LeagueProcessor
}

// *****************************************************************************
// public
// *****************************************************************************

// NewEventProcessor creates a new struct for scanning and storing crypto soccer events
func NewEventProcessor(
	client *ethclient.Client,
	db *storage.Storage,
	engine *engine.Engine,
	assets *assets.Assets,
	leagues *leagues.Leagues,
	updates *updates.Updates,
	market *market.Market,
) (*EventProcessor, error) {
	divisionCreationProcessor, err := NewDivisionCreationProcessor(db, assets, leagues)
	if err != nil {
		return nil, err
	}
	leagueProcessor, err := NewLeagueProcessor(engine, leagues, db)
	if err != nil {
		return nil, err
	}
	return &EventProcessor{
		client,
		db,
		engine,
		assets,
		leagues,
		updates,
		market,
		divisionCreationProcessor,
		leagueProcessor,
	}, nil
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

	scanner := NewEventScanner(p.assets, p.updates, p.market)
	if scanner == nil {
		return opts.Start, errors.New("Unable to create scanner")
	}
	if err := scanner.Process(opts); err != nil {
		return 0, err
	} else {
		for _, v := range scanner.Events {
			if err := p.dispatch(v); err != nil {
				return 0, err
			}
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
	log.Debugf("[process] dispach event block %v inBlockIndex %v", e.BlockNumber, e.TxIndexInBlock)

	switch v := e.Value.(type) {
	case assets.AssetsDivisionCreation:
		log.Debug("[processor] Dispatching LeaguesDivisionCreation event")
		return p.divisionCreationProcessor.Process(v)
	case assets.AssetsTeamTransfer:
		log.Debug("[processor] dispatching LeaguesTeamTransfer event")
		teamID := v.TeamId
		newOwner := v.To.String()
		team, err := p.db.GetTeam(teamID)
		if err != nil {
			return err
		}
		// team.State.BlockNumber = blockNumber
		team.State.Owner = newOwner
		return p.db.TeamUpdate(teamID, team.State)
	case assets.AssetsPlayerStateChange:
		log.Debug("[processor] dispatching LeaguesPlayerStateChange event")
		playerID := v.PlayerId
		state := v.State
		shirtNumber, err := p.assets.GetCurrentShirtNum(&bind.CallOpts{}, state)
		if err != nil {
			return err
		}
		teamID, err := p.assets.GetCurrentTeamId(&bind.CallOpts{}, state)
		if err != nil {
			return err
		}
		player, err := p.db.GetPlayer(playerID)
		if err != nil {
			return err
		}
		player.State.TeamId = teamID
		player.State.ShirtNumber = uint8(shirtNumber.Uint64())
		return p.db.PlayerUpdate(playerID, player.State)
	case updates.UpdatesActionsSubmission:
		log.Debug("[processor] Dispatching UpdatesActionsSubmission event")
		return p.leagueProcessor.Process(v)
	case market.MarketPlayerFreeze:
		log.Debug("[processor] Dispatching MarketPlayerFreeze event")
		playerID := v.PlayerId
		player, err := p.db.GetPlayer(playerID)
		if err != nil {
			return err
		}
		player.State.Frozen = v.Frozen
		return p.db.PlayerUpdate(playerID, player.State)
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
	block, err := p.client.BlockByHash(context.Background(), eventRaw.BlockHash)
	if err != nil {
		return 0, 0, err
	}
	return block.Time(), eventRaw.BlockNumber, nil
}
