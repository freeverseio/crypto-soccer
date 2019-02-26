package service

import (
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	cfg "github.com/freeverseio/go-soccer/config"
	"github.com/freeverseio/go-soccer/eth"

	sto "github.com/freeverseio/go-soccer/storage"
	log "github.com/sirupsen/logrus"
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
	lionelAddress := common.HexToAddress(cfg.C.Contracts.LionelAddress)

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
	log.Info("Waiting terminate channel")
	<-s.terminatedch
}

func (s *Service) updateLeague(leagueNo int64) error {

	var hash common.Hash
	tx, _, err := s.lionel.SendTransactionSync(nil, 0,
		"update",
		big.NewInt(leagueNo), &hash)

	if err == nil {
		log.WithField("tx", tx.Hash().Hex()).Info("  League ", leagueNo, " : updating")
	} else {
		log.Error("  League ", leagueNo, " : update failed")
	}
	return err
}

func (s *Service) challangeLeague(leagueNo int64) error {

	var hash common.Hash
	tx, _, err := s.lionel.SendTransactionSync(nil, 0,
		"challange",
		big.NewInt(leagueNo), hash)

	if err == nil {
		log.WithField("tx", tx.Hash().Hex()).Info("  League ", leagueNo, " : challanging")
	} else {
		log.Error("  League ", leagueNo, " : challange failed")
	}
	return err
}

func (s *Service) process() (bool, error) {

	var legueCount *big.Int
	if err := s.lionel.Call(&legueCount, "legueCount"); err != nil {
		return false, err
	}
	log.Info("Scanning ", legueCount.Uint64(), " leagues...")

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
		log.Info("Starting service...")
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
