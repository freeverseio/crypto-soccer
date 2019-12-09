package process

import (
	"time"

	"github.com/freeverseio/crypto-soccer/go/contracts"
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
	contracts *contracts.Contracts,
	universedb *storage.Storage,
	relaydb *relay.Storage,
	namesdb *names.Generator,
) (*BackgroundProcess, error) {
	eventProcessor, err := NewEventProcessor(
		contracts,
		universedb,
		relaydb,
		namesdb,
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

func (b *BackgroundProcess) Process() (uint64, error) {
	err := b.eventProcessor.universedb.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			b.eventProcessor.universedb.Rollback()
			return
		}
		b.eventProcessor.universedb.Commit()
	}()
	err = b.eventProcessor.relaydb.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			b.eventProcessor.relaydb.Rollback()
			return
		}
		b.eventProcessor.relaydb.Commit()
	}()

	delta := uint64(1000)
	return b.eventProcessor.Process(delta)
}

func (b *BackgroundProcess) Start() {
	go func() {
	L:
		for {
			select {
			case <-b.queryStop:
				break L
			default:
				processedBlocks, err := b.Process()
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
