package process

import (
	"time"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
)

type BackgroundProcess struct {
	assetsContract *assets.Assets
	storage        *storage.Storage
	queryStop      chan (bool)
	stopped        chan (bool)
}

func BackgroundProcessNew(assetsContract *assets.Assets, storage *storage.Storage) *BackgroundProcess {
	return &BackgroundProcess{
		assetsContract: assetsContract,
		storage:        storage,
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
				Process(b.assetsContract, b.storage)
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
