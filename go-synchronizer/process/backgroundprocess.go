package process

import (
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/states"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	log "github.com/sirupsen/logrus"
)

type BackgroundProcess struct {
	eventProcessor *EventProcessor
	queryStop      chan (bool)
	stopped        chan (bool)
}

func BackgroundProcessNew(
	client *ethclient.Client,
	assetsContract *assets.Assets,
	statesContract *states.States,
	leaguesContract *leagues.Leagues,
	storage *storage.Storage,
) *BackgroundProcess {
	return &BackgroundProcess{
		eventProcessor: NewEventProcessor(client, storage, assetsContract, statesContract, leaguesContract),
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
					log.Error(err)
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
