package process

import (
	"context"
	"errors"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/updates"

	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	//"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	log "github.com/sirupsen/logrus"
)

type EventProcessor struct {
	usesGanache               bool
	client                    *ethclient.Client
	db                        *storage.Storage
	engine                    *engine.Engine
	leagues                   *leagues.Leagues
	updates                   *updates.Updates
	divisionCreationProcessor *DivisionCreationProcessor
	leagueProcessor           *LeagueProcessor
}

// *****************************************************************************
// public
// *****************************************************************************

// NewEventProcessor creates a new struct for scanning and storing crypto soccer events
func NewEventProcessor(client *ethclient.Client, db *storage.Storage, engine *engine.Engine, leagues *leagues.Leagues, updates *updates.Updates) (*EventProcessor, error) {
	divisionCreationProcessor, err := NewDivisionCreationProcessor(db, leagues)
	if err != nil {
		return nil, err
	}
	leagueProcessor, err := NewLeagueProcessor(engine, leagues, db)
	if err != nil {
		return nil, err
	}
	return &EventProcessor{
		false,
		client,
		db,
		engine,
		leagues,
		updates,
		divisionCreationProcessor,
		leagueProcessor,
	}, nil
}

// NewGanacheEventProcessor creates a new struct for scanning and storing crypto soccer events from a ganache client
func NewGanacheEventProcessor(client *ethclient.Client, db *storage.Storage, engine *engine.Engine, leagues *leagues.Leagues, updates *updates.Updates) (*EventProcessor, error) {
	divisionCreationProcessor, err := NewDivisionCreationProcessor(db, leagues)
	if err != nil {
		return nil, err
	}
	leagueProcessor, err := NewLeagueProcessor(engine, leagues, db)
	if err != nil {
		return nil, err
	}
	return &EventProcessor{true, client, db, engine, leagues, updates, divisionCreationProcessor, leagueProcessor}, nil
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

	scanner := NewEventScanner()
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
	if err := scanner.Process(); err != nil {
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
	log.Debugf("[process] dispach event block %v inBlockIndex %v", e.BlockNumber, e.TxIndexInBlock)

	switch v := e.Value.(type) {
	case leagues.LeaguesDivisionCreation:
		log.Debug("[processor] Dispatching LeaguesDivisionCreation event")
		return p.divisionCreationProcessor.Process(v)
	case leagues.LeaguesTeamTransfer:
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
	case leagues.LeaguesPlayerStateChange:
		log.Debug("[processor] dispatching LeaguesPlayerStateChange event")
		playerID := v.PlayerId
		state := v.State
		shirtNumber, err := p.leagues.GetCurrentShirtNum(&bind.CallOpts{}, state)
		if err != nil {
			return err
		}
		teamID, err := p.leagues.GetCurrentTeamId(&bind.CallOpts{}, state)
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
