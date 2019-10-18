package relay

import (
	"crypto/ecdsa"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/freeverseio/crypto-soccer/relay/contracts/updates"
	"github.com/freeverseio/crypto-soccer/relay/storage"
)

type BackgroundProcess struct {
	relay     *Processor
	queryStop chan (bool)
	stopped   chan (bool)
}

func BackgroundProcessNew(
	client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	storage *storage.Storage,
	updatesContract *updates.Updates,
) (*BackgroundProcess, error) {
	processor, err := NewProcessor(client, privateKey, storage, updatesContract)
	if err != nil {
		return nil, err
	}
	return &BackgroundProcess{
		relay:     processor,
		queryStop: make(chan (bool)),
		stopped:   make(chan (bool)),
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
				err := b.relay.Process()
				if err != nil {
					panic(err)
				}
				time.Sleep(1 * time.Second)
			}
		}
		b.stopped <- true
	}()
}

func (b *BackgroundProcess) StopAndJoin() {
	b.queryStop <- true
	<-b.stopped
}
