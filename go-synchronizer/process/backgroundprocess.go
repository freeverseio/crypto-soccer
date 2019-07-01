package process

import (
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
)

type BackgroundProcess struct {
	client         *ethclient.Client
	assetsContract *assets.Assets
	storage        *storage.Storage
	eventProcessor *EventProcessor
	queryStop      chan (bool)
	stopped        chan (bool)
}

func BackgroundProcessNew(client *ethclient.Client, assetsContract *assets.Assets, storage *storage.Storage) *BackgroundProcess {
	return &BackgroundProcess{
		client:         client,
		assetsContract: assetsContract,
		storage:        storage,
		eventProcessor: NewEventProcessor(nil, storage, assetsContract),
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
				b.eventProcessor.Process()
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
