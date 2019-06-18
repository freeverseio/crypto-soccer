package process

import (
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
)

type BackgroundProcess struct {
	assetsContract *assets.Assets
	sto            storage.Storage
	queryStop      chan (bool)
	stopped        chan (bool)
}

func BackgroundProcessNew() *BackgroundProcess {
	return &BackgroundProcess{
		queryStop: make(chan (bool)),
		stopped:   make(chan (bool)),
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
				// log.Info("tick")
			}
		}
		b.stopped <- true
	}()
}

func (b *BackgroundProcess) StopAndJoin() {
	b.queryStop <- true
	<-b.stopped
}
