package process

import (
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"
)

type BackgroundProcess struct {
	assetsContract *assets.Assets
	sto storage.Storage
	queryStop chan(bool)
}

// func New(assetsContract *assets.Assets,sto storage.Storage) *BackgroundProcess {
// 	return &BackgroundProcess{
// 		assetsContract,
// 		sto,
// 		queryStop: new(chan(bool)),
// 		stopped: new(chan(bool)),
// 	}
// }

// func (b *BackgroundProcess) start() {
// 	go func(){
// 		for {
// 			select {
// 				case <-b.queryStop : break
// 				default:
// 					Process(b.assetsContract,b.sto)
// 					opcount++
// 			}				
// 		}
// 		stopped <- true	
// 	}()
// }

// func (b *BackgroundProcess) stopAndJoin() {
// 	b.queryStop <- true
// 	<- b.stopped
// }