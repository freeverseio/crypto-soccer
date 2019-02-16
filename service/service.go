package service

import (
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	cfg "github.com/freeverseio/go-soccer/config"
	"github.com/freeverseio/go-soccer/eth"

	sto "github.com/freeverseio/go-soccer/storage"
)

type Service struct {
	storage *sto.Storage
	web3    *eth.Web3Client
	lionel  *eth.Contract

	stats     ServiceStats
	laststats ServiceStats

	terminatech  chan interface{}
	terminatedch chan interface{}
}

var (
	errVerifySmartcontract = errors.New("cannot verify deployed smartcontract")
	errReadPersistLimit    = errors.New("error reading current persistLimit")
	errReachedPersistLimit = errors.New("persistlimit reached")
)

func NewService(web3 *eth.Web3Client, storage *sto.Storage) (*Service, error) {

	lionelAbi, err := abi.JSON(strings.NewReader(lionelAbiJson))
	if err != nil {
		return nil, err
	}
	lionelAddress := common.HexToAddress(cfg.C.LionelAddress)

	lionel, err := eth.NewContract(web3, &lionelAbi, nil, &lionelAddress)
	if err != nil {
		return nil, err
	}

	return &Service{
		web3:         web3,
		storage:      storage,
		lionel:       lionel,
		terminatech:  make(chan interface{}),
		terminatedch: make(chan interface{}),
	}, nil
}

// Stop scanning the blockchain for events
func (s *Service) Stop() {
	go func() {
		s.terminatech <- nil
	}()
}

// Join waits all background jobs finished
func (s *Service) Join() {
	<-s.terminatedch
}

func (s *Service) updateLeague(leagueNo int64) error {
	return nil
}

func (s *Service) challangeLeague(leagueNo int64) error {
	return nil
}

func (s *Service) process() (bool, error) {

	var legueCount *big.Int
	if err := s.lionel.Call(&legueCount, "legueCount"); err != nil {
		return false, err
	}

	for i := int64(0); i < int64(legueCount.Uint64()); i++ {

		legueNo := big.NewInt(i)

		var canLeagueBeUpdated bool
		if err := s.lionel.Call(&canLeagueBeUpdated, "canLeagueBeUpdated", legueNo); err != nil {
			return false, err
		}
		if canLeagueBeUpdated {
			if err := s.updateLeague(i); err != nil {
				return false, err
			}
		}

		var canLeagueBeChallanged bool
		if err := s.lionel.Call(&canLeagueBeChallanged, "canLeagueBeChallanged", legueNo); err != nil {
			return false, err
		}
		if canLeagueBeChallanged {
			if err := s.challangeLeague(i); err != nil {
				return false, err
			}
		}

	}

	return false, nil
}

// Start scanning the blockchain for events
func (s *Service) Start() {

	go func() {
		loop := true
		for loop {
			select {

			case <-s.terminatech:
				log.Debug("EVENT Dispatching terminatech")
				loop = false
				break

			default:
				finished, err := s.process()
				if err != nil {
					log.Error("EVENT Failed ", err)
					loop = false
				} else if finished {
					time.Sleep(4 * time.Second)
				}
			}
		}
		s.terminatedch <- nil
	}()
}
