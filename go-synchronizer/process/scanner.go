package process

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/updates"
)

type AbstractEvent struct {
	BlockNumber    uint64
	TxIndexInBlock uint
	Name           string // for ease of debugging
	Value          interface{}
}

func NewAbstractEvent(log types.Log, name string, x interface{}) *AbstractEvent {
	return &AbstractEvent{log.BlockNumber, log.TxIndex, name, x}
}

type EventScanner struct {
	leagues *leagues.Leagues
	updates *updates.Updates
	Events  []*AbstractEvent
}

func NewEventScanner(leagues *leagues.Leagues, updates *updates.Updates) *EventScanner {
	return &EventScanner{leagues, updates, []*AbstractEvent{}}
}

// Process: scans all events types and puts them in the Events slice
// types of events it listens
// leagues.LeaguesDivisionCreation
// leagues.LeaguesTeamTransfer
// leagues.LeaguesPlayerTransfer
// updates.UpdatesActionsSubmission

func (s *EventScanner) Process(opts *bind.FilterOpts) error {
	if err := s.scanDivisionCreation(opts); err != nil {
		return err
	}
	if err := s.scanTeamTransfer(opts); err != nil {
		return err
	}
	if err := s.scanPlayerTransfer(opts); err != nil {
		return err
	}
	if err := s.scanActionsSubmission(opts); err != nil {
		return err
	}

	// TODO: sort all events

	return nil
}

// TODO: abstract away all scan functions by passing the Filterfunction and name
func (s *EventScanner) scanDivisionCreation(opts *bind.FilterOpts) error {
	if opts == nil {
		opts = &bind.FilterOpts{Start: 0}
	}
	iter, err := s.leagues.FilterDivisionCreation(opts)
	if err != nil {
		return err
	}

	for iter.Next() {
		e := *(iter.Event)
		s.Events = append(s.Events, NewAbstractEvent(e.Raw, "LeaguesDivisionCreation", e))
	}
	return nil
}

func (s *EventScanner) scanTeamTransfer(opts *bind.FilterOpts) error {
	if opts == nil {
		opts = &bind.FilterOpts{Start: 0}
	}
	iter, err := s.leagues.FilterTeamTransfer(opts)
	if err != nil {
		return err
	}

	for iter.Next() {
		e := *(iter.Event)
		s.Events = append(s.Events, NewAbstractEvent(e.Raw, "LeaguesTeamTransfer", e))
		log.Debug("Adding event in scanner.scanTeamTransfer")
	}
	return nil
}

func (s *EventScanner) scanPlayerTransfer(opts *bind.FilterOpts) error {
	if opts == nil {
		opts = &bind.FilterOpts{Start: 0}
	}
	iter, err := s.leagues.FilterPlayerTransfer(opts)
	if err != nil {
		return err
	}

	for iter.Next() {
		e := *(iter.Event)
		s.Events = append(s.Events, NewAbstractEvent(e.Raw, "LeaguesPlayerTransfer", e))
		log.Debug("Adding event in scanner.scanPlayerTransfer")
	}
	return nil
}

func (s *EventScanner) scanActionsSubmission(opts *bind.FilterOpts) error {
	if opts == nil {
		opts = &bind.FilterOpts{Start: 0}
	}
	iter, err := s.updates.FilterActionsSubmission(opts)
	if err != nil {
		return err
	}

	for iter.Next() {
		e := *(iter.Event)
		s.Events = append(s.Events, NewAbstractEvent(e.Raw, "UpdatesActionsSubmission", e))
		log.Debug("Adding event in scanner.scanActionsSubmission")
	}
	return nil
}
