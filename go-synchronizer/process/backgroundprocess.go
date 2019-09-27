package process

import (
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"
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
	marketContract *market.Market,
	assetsContract *assets.Assets,
	updatesContract *updates.Updates,
) *BackgroundProcess {
	return &BackgroundProcess{
		eventProcessor: NewEventProcessor(client, storage, marketContract, assetsContract, updatesContract),
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
