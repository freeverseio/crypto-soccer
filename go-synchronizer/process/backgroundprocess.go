package process

import (
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
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
) *BackgroundProcess {
	return &BackgroundProcess{
		eventProcessor: NewEventProcessor(client, storage, engineContract, leaguesContract, updatesContract),
		queryStop:      make(chan (bool)),
		stopped:        make(chan (bool)),
	}
}

func (b *BackgroundProcess) Start() {
	go func() {
	L:
		for {
			select {
			case <-b.queryStop:
				break L
			default:
				err := b.eventProcessor.Process()
				if err != nil {
					panic(err)
				}
				time.Sleep(2 * time.Second)
			}
		}
		b.stopped <- true
	}()
}

func (b *BackgroundProcess) StopAndJoin() {
	b.queryStop <- true
	<-b.stopped
}
