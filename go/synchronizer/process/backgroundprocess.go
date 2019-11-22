package process

import (
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/freeverseio/crypto-soccer/go/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go/contracts/engineprecomp"
	"github.com/freeverseio/crypto-soccer/go/contracts/evolution"
	"github.com/freeverseio/crypto-soccer/go/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go/contracts/market"
	"github.com/freeverseio/crypto-soccer/go/contracts/updates"
	"github.com/freeverseio/crypto-soccer/go/names"
	relay "github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

type BackgroundProcess struct {
	eventProcessor *EventProcessor
	queryStop      chan (bool)
	stopped        chan (bool)
}

func BackgroundProcessNew(
	client *ethclient.Client,
	universedb *storage.Storage,
	relaydb *relay.Storage,
	namesdb *names.Generator,
	engineContract *engine.Engine,
	enginePreCompContract *engineprecomp.Engineprecomp,
	assetsContract *assets.Assets,
	leaguesContract *leagues.Leagues,
	updatesContract *updates.Updates,
	marketContract *market.Market,
	evolutionContract *evolution.Evolution,
) (*BackgroundProcess, error) {
	eventProcessor, err := NewEventProcessor(
		client,
		universedb,
		relaydb,
		namesdb,
		engineContract,
		enginePreCompContract,
		assetsContract,
		leaguesContract,
		updatesContract,
		marketContract,
		evolutionContract,
	)

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
