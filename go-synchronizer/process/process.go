package process

import (
	"context"
	//"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/states"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	log "github.com/sirupsen/logrus"
)

type EventProcessor struct {
	client *ethclient.Client
	db     *storage.Storage
	assets *assets.Assets
	states *states.States
}

// *****************************************************************************
// public
// *****************************************************************************

// NewEventProcessor creates a new struct for scanning and storing crypto soccer events
func NewEventProcessor(client *ethclient.Client, db *storage.Storage, assets *assets.Assets, states *states.States) *EventProcessor {
	return &EventProcessor{client, db, assets, states}
}

// Process processes all scanned events and stores them into the database db
func (p *EventProcessor) Process() error {
	opts := p.nextRange()

	if opts == nil {
		log.Info("No new blocks to search for events")
		return nil
	}

	log.WithFields(log.Fields{
		"start": opts.Start,
		"end":   *opts.End,
	}).Info("Syncing ...")

	// scan TeamCreated events in range [start, end]
	if events, err := p.scanTeamCreated(opts); err != nil {
		return err
	} else {
		err = p.storeTeamCreated(events)
		if err != nil {
			return err
		}
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
func (p *EventProcessor) storeTeamCreated(events []assets.AssetsTeamCreated) error {
	for _, event := range events {
		if name, err := p.assets.GetTeamName(nil, event.Id); err != nil {
			return err
		} else if err := p.db.TeamAdd(&storage.Team{event.Id.Uint64(), name}); err != nil {
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
func (p *EventProcessor) storeVirtualPlayers(teamId *big.Int) error {
	for i := 0; i < 11; i++ {
		if id, err := p.assets.GenerateVirtualPlayerId(&bind.CallOpts{}, teamId, uint8(i)); err != nil {
			return err
		} else if state, err := p.assets.GenerateVirtualPlayerState(&bind.CallOpts{}, id); err != nil {
			return err
		} else {
			if skills, err := p.states.GetSkillsVec(&bind.CallOpts{}, state); err != nil {
				return err
			} else {
				player := storage.Player{
					Id:        id.Uint64(),
					TeamId:    teamId.Uint64(),
					State:     state.String(),
					Defence:   uint64(skills[0]),
					Speed:     uint64(skills[1]),
					Pass:      uint64(skills[2]),
					Shoot:     uint64(skills[3]),
					Endurance: uint64(skills[4]),
				}
				p.db.PlayerAdd(&player)
				if stored, err := p.db.GetPlayer(id.Uint64()); err != nil {
					log.Fatal(err)
				} else if stored.State != state.String() {
					log.Fatal("Mismatch while storing virtual player. State before storage:", state.String(), " vs state after storage:", stored.Id, stored.State)
				}
			}
		}
	}
	return nil
}
