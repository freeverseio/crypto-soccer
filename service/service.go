package service

import (
	"errors"
	"time"

	"github.com/freeverseio/go-soccer/eth"
	"github.com/freeverseio/go-soccer/lionel"

	sto "github.com/freeverseio/go-soccer/storage"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	storage *sto.Storage
	web3    *eth.Web3Client
	lionel  *lionel.Lionel
	stakers *eth.Contract

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

	// load lionel
	lionel, err := lionel.New(web3, storage)
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
				lionelFinished, err := s.processLionel()
				if err != nil {
					log.Error("Lionel failed ", err)
				}
				stakersFinished, err := s.processStakers()
				if err != nil {
					log.Error("Stakers failed ", err)
				}
				if lionelFinished && stakersFinished {
					log.Info("All finished, taking a litte nap ")
					time.Sleep(4 * time.Second)
				}
			}
		}
		s.terminatedch <- nil
	}()
}
