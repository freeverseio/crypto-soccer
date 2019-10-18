package process

import (
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/market"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/updates"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
)

type BackgroundProcess struct {
	eventProcessor *EventProcessor
	queryStop      chan (bool)
	stopped        chan (bool)
}

func BackgroundProcessNew(
	client *ethclient.Client,
	storage *storage.Storage,
	engineContract *engine.Engine,
	leaguesContract *leagues.Leagues,
	updatesContract *updates.Updates,
	marketContract *market.Market,
) (*BackgroundProcess, error) {
	eventProcessor, err := NewEventProcessor(client, storage, engineContract, leaguesContract, updatesContract, marketContract)
	if err != nil {
		return nil, err
	}
	return &BackgroundProcess{
		eventProcessor: eventProcessor,
		queryStop:      make(chan (bool)),
		stopped:        make(chan (bool)),
	}, nil
}

func (b *BackgroundProcess) Start() {
	go func() {
	L:
		for {
			select {
			case <-b.queryStop:
				break L
			default:
				delta := uint64(1000)
				processedBlocks, err := b.eventProcessor.Process(delta)
				if err != nil {
					panic(err)
				}
				if processedBlocks == 0 {
					time.Sleep(2 * time.Second)
				}
			}
		}
		b.stopped <- true
	}()
}

func (b *BackgroundProcess) StopAndJoin() {
	b.queryStop <- true
	<-b.stopped
}
