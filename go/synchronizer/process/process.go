package process

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/freeverseio/crypto-soccer/go/useractions"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/contracts/autogenerated/assets"
	"github.com/freeverseio/crypto-soccer/go/contracts/autogenerated/proxy"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/freeverseio/crypto-soccer/go/contracts/autogenerated/market"
	"github.com/freeverseio/crypto-soccer/go/contracts/autogenerated/updates"
	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/staker"
	log "github.com/sirupsen/logrus"
)

type EventProcessor struct {
	Client               *ethclient.Client
	contracts            *contracts.Contracts
	proxyContractAddress string
	staker               *staker.Staker
	namesdb              *names.Generator
	useractionsPublisher useractions.UserActionsPublishService
	blocksUntilFinal     uint64
}

// *****************************************************************************
// public
// *****************************************************************************

// NewEventProcessor creates a new struct for scanning and storing crypto soccer events
func NewEventProcessor(
	client *ethclient.Client,
	proxyContractAddress string,
	namesdb *names.Generator,
	useractionsPublisher useractions.UserActionsPublishService,
	staker *staker.Staker,
	blocksUntilFinal uint64,
) *EventProcessor {
	return &EventProcessor{
		client,
		nil,
		proxyContractAddress,
		staker,
		namesdb,
		useractionsPublisher,
		blocksUntilFinal,
	}
}

func (b EventProcessor) getFirstDeployEvent(proxyAddress string) (*proxy.ProxyNewDirectory, error) {
	proxyContract, err := proxy.NewProxy(common.HexToAddress(proxyAddress), b.Client)
	if err != nil {
		return nil, err
	}

	iter, err := proxyContract.FilterNewDirectory(&bind.FilterOpts{
		Start:   0,
		End:     nil,
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}
	if !iter.Next() {
		return nil, nil
	}

	return iter.Event, nil
}

func (p *EventProcessor) Process(tx *sql.Tx, delta uint64) (uint64, error) {
	lastProcessedBlock, err := p.dbLastBlockNumber(tx)
	if err != nil {
		return 0, err
	}
	if lastProcessedBlock == 0 {
		bigBangEvent, err := p.getFirstDeployEvent(p.proxyContractAddress)
		if err != nil {
			return 0, err
		}
		if bigBangEvent == nil {
			return 0, err
		}

		log.Infof("[processor|consume] bigBang event ... block %v , directory %v", bigBangEvent.Raw.BlockNumber, bigBangEvent.Addr.Hex())

		c, err := contracts.NewByNewDirectoryEvent(p.Client, *bigBangEvent)
		if err != nil {
			return 0, err
		}
		if p.staker != nil {
			if err := p.staker.Init(*c); err != nil {
				return 0, err
			}
		}
		bigBangBlock := bigBangEvent.Raw.BlockNumber
		if err = storage.SetBlockNumber(tx, bigBangBlock); err != nil {
			return 0, err
		}
		// TODO save contacts to postgres
		if err := c.ToStorage(tx); err != nil {
			return 0, err
		}
		p.contracts = c
		return bigBangBlock, nil
	}

	if p.contracts, err = contracts.NewFromStorage(p.Client, tx); err != nil {
		return 0, err
	}
	return p.Process2(tx, delta)
}

// Process processes all scanned events and stores them into the database db
func (p *EventProcessor) Process2(tx *sql.Tx, delta uint64) (uint64, error) {
	opts, err := p.nextRange(tx, 1)
	if err != nil {
		return 0, err
	}

	if opts == nil {
		log.Debug("No new blocks to scan.")
		return 0, nil
	}

	log.WithFields(log.Fields{
		"start": opts.Start,
		"end":   *opts.End,
	}).Info("digest ...")

	scanner := NewEventScanner(p.contracts)
	if scanner == nil {
		return opts.Start, errors.New("Unable to create scanner")
	}

	if err := scanner.Process(opts); err != nil {
		return 0, err
	}

	for _, v := range scanner.Events {
		if err := p.Dispatch(tx, v); err != nil {
			return 0, err
		}
	}

	err = storage.SetBlockNumber(tx, *opts.End)
	deltaBlock := *opts.End - opts.Start + 1

	return deltaBlock, err
}

// *****************************************************************************
// private
// *****************************************************************************
func (p *EventProcessor) Dispatch(tx *sql.Tx, e *AbstractEvent) error {
	log.Debugf("[process] dispach event block %v inBlockIndex %v", e.BlockNumber, e.TxIndexInBlock)

	switch v := e.Value.(type) {
	case proxy.ProxyNewDirectory:
		var err error
		p.contracts, err = ConsumeNewDirectory(tx, *p.contracts, v)
		return err
	case assets.AssetsAssetsInit:
		assetsInitProcessor := NewAssetsInitProcessor(p.contracts)
		return assetsInitProcessor.Process(tx, v)
	case assets.AssetsDivisionCreation:
		divisionCreationProcessor := NewDivisionCreationProcessor(p.contracts, p.namesdb)
		return divisionCreationProcessor.Process(tx, v)
	case assets.AssetsTeamTransfer:
		return ConsumeTeamTransfer(tx, v)
	case market.MarketPlayerStateChange:
		return ConsumePlayerStateChange(tx, p.contracts, p.namesdb, v)
	case updates.UpdatesActionsSubmission:
		return ConsumeActionsSubmission(tx, p.contracts, p.useractionsPublisher, p.staker, v)
	case updates.UpdatesTimeZoneUpdate:
		return ConsumeTimezoneUpdate(tx, v)
	default:
		return fmt.Errorf("[processor|consume] unknown event %+v", e)
	}
}
func (p *EventProcessor) nextRange(tx *sql.Tx, delta uint64) (*bind.FilterOpts, error) {
	start, err := p.dbLastBlockNumber(tx)
	if err != nil {
		return nil, err
	}
	if start != 0 {
		// unless this is the very first execution,
		// the block number that is stored in the db
		// was already scanned. We are interested in
		// the next block
		if start < math.MaxUint64 {
			start += 1
		} else {
			return nil, errors.New("Block range overflow")
		}
	}
	end := p.clientLastBlockNumber() - p.blocksUntilFinal
	if delta != 0 {
		end = uint64(math.Min(float64(start+delta-1), float64(end)))
	}
	if start > end {
		return nil, nil
	}
	return &bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	}, nil
}

func (p *EventProcessor) clientLastBlockNumber() uint64 {
	if p.contracts.Client == nil {
		log.Warn("Client is nil. Returning 0 as last block.")
		return 0
	}
	header, err := p.contracts.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Warn("Could not get blockchain last block")
		return 0
	}
	return header.Number.Uint64()
}
func (p *EventProcessor) dbLastBlockNumber(tx *sql.Tx) (uint64, error) {
	storedLastBlockNumber, err := storage.GetBlockNumber(tx)
	if err != nil {
		return 0, err
	}
	return storedLastBlockNumber, err
}
func (p *EventProcessor) getTimeOfEvent(eventRaw types.Log) (uint64, uint64, error) {
	block, err := p.contracts.Client.BlockByHash(context.Background(), eventRaw.BlockHash)
	if err != nil {
		return 0, 0, err
	}
	return block.Time(), eventRaw.BlockNumber, nil
}
