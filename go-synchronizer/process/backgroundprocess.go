package process

import (
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"
)

type BackgroundProcess struct {
	assetsContract *assets.Assets
	sto storage.Storage
	queryStop chan(bool)
	stopped chan(bool)
}

func BackgroundProcessNew() *BackgroundProcess {
	return &BackgroundProcess{
		queryStop: make(chan(bool)),
		stopped: make(chan(bool)),
	}
}

func (b *BackgroundProcess) Start() {
	go func(){
		for {
			select {
				case <-b.queryStop : break
				default: Process(b.assetsContract,b.sto)
			}				
			b.stopped <- true	
		}
	}()
}

func (b *BackgroundProcess) StopAndJoin() {
	b.queryStop <- true
	<- b.stopped
}
