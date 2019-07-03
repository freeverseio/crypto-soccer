package process

import (
	"context"
	//"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	log "github.com/sirupsen/logrus"
)

type EventProcessor struct {
	client *ethclient.Client
	db     *storage.Storage
	assets *assets.Assets
}

// *****************************************************************************
// public
// *****************************************************************************

// NewEventProcessor creates a new struct for scanning and storing crypto soccer events
func NewEventProcessor(client *ethclient.Client, db *storage.Storage, assets *assets.Assets) *EventProcessor {
	return &EventProcessor{client, db, assets}
}

// Process processes all scanned events and stores them into the database db
func (p *EventProcessor) Process() error {
	log.Info("Syncing ...")
	log.Trace("Process: scanning the blockchain")

	start, end := p.nextRange()

	if start > end {
		log.Debug("No new blocks to search for events")
		return nil
	}

	opts := &bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	}

	// scan TeamCreated events in range [start, end]
	if events, err := p.scanTeamCreated(opts); err != nil {
		return err
	} else {
		p.storeTeamCreated(events)
	}

	// update the store block in the database
	p.db.SetBlockNumber(big.NewInt(int64(end + 1)))
	return nil
}

// *****************************************************************************
// private
// *****************************************************************************

func (p *EventProcessor) nextRange() (uint64, uint64) {
	return p.dbLastBlockNumber(), p.clientLastBlockNumber()
}
func (p *EventProcessor) clientLastBlockNumber() uint64 {
	if p.client == nil {
		return 0
	}
	header, err := p.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Println(err)
		return 0
	}
	return header.Number.Uint64()
}
func (p *EventProcessor) dbLastBlockNumber() uint64 {
	storedLastBlockNumber, err := p.db.GetBlockNumber()
	if err != nil {
		log.Println(err)
		return 0
	}
	return storedLastBlockNumber.Uint64()
}
func (p *EventProcessor) storeTeamCreated(events []assets.AssetsTeamCreated) error {
	for _, event := range events {
		if name, err := p.assets.GetTeamName(nil, event.Id); err != nil {
			return err
		} else if err := p.db.TeamAdd(event.Id.Uint64(), name); err != nil {
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
			p.db.PlayerAdd(&storage.Player{id.Uint64(), state})
			log.Println("virtual player before storage:", id.Uint64(), state.String())
			if stored, err := p.db.GetPlayer(id.Uint64()); err != nil {
				log.Fatal(err)
			} else {
				log.Println("virtual player stored:", stored.Id, stored.State.String())
			}
		}
	}
	return nil
}
