package process

import (
	"database/sql"
	"time"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

type BackgroundProcess struct {
	eventProcessor *EventProcessor
	queryStop      chan (bool)
	stopped        chan (bool)
	relaydb        *sql.DB
}

func BackgroundProcessNew(
	contracts *contracts.Contracts,
	universedb *storage.Storage,
	relaydb *sql.DB,
	namesdb *names.Generator,
) (*BackgroundProcess, error) {
	eventProcessor, err := NewEventProcessor(
		contracts,
		universedb,
		namesdb,
	)

	if err != nil {
		return nil, err
	}
	return &BackgroundProcess{
		eventProcessor: eventProcessor,
		queryStop:      make(chan (bool)),
		stopped:        make(chan (bool)),
		relaydb:        relaydb,
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
	tx, err := b.relaydb.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
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
