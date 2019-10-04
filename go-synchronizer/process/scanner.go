package process

import (
	"sort"

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

func NewAbstractEvent(blockNumber uint64, txIndex uint, name string, x interface{}) *AbstractEvent {
	return &AbstractEvent{blockNumber, txIndex, name, x}
}

type EventScanner struct {
	leagues *leagues.Leagues
	updates *updates.Updates
	Events  []*AbstractEvent
}

func NewEventScanner(leagues *leagues.Leagues, updates *updates.Updates) *EventScanner {
	return &EventScanner{leagues, updates, []*AbstractEvent{}}
}

type byFunction func(p1, p2 *AbstractEvent) bool

func (b byFunction) Sort(events []*AbstractEvent) {
	ps := &abstractEventSorter{
		events: events,
		by:     b,
	}
	sort.Sort(ps)
}

type abstractEventSorter struct {
	events []*AbstractEvent
	by     func(p1, p2 *AbstractEvent) bool // Closure used in the Less method.
}

func (s *abstractEventSorter) Len() int {
	return len(s.events)
}

func (s *abstractEventSorter) Swap(i, j int) {
	s.events[i], s.events[j] = s.events[j], s.events[i]
}

func (s *abstractEventSorter) Less(i, j int) bool {
	return s.by(s.events[i], s.events[j])
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

	if len(s.Events) > 10000 {
		log.Info("sorting ", len(s.Events), " be patient...")
	}

	sortFunction := func(p1, p2 *AbstractEvent) bool {
		if p1.BlockNumber == p2.BlockNumber {
			return p1.TxIndexInBlock < p2.TxIndexInBlock
		}
		return p1.BlockNumber < p2.BlockNumber

	}

	byFunction(sortFunction).Sort(s.Events)

	return nil
}

func (s *EventScanner) addEvent(rawEvent types.Log, name string, event interface{}) {
	log.Info("Add event ", name)
	s.Events = append(s.Events, NewAbstractEvent(rawEvent.BlockNumber, rawEvent.TxIndex, name, event))
}

func (s *EventScanner) scanDivisionCreation(opts *bind.FilterOpts) error {
	iter, err := s.leagues.FilterDivisionCreation(opts)
	if err != nil {
		return err
	}

	for iter.Next() {
		e := *(iter.Event)
		s.addEvent(e.Raw, "LeaguesDivisionCreation", e)
	}
	return nil
}

func (s *EventScanner) scanTeamTransfer(opts *bind.FilterOpts) error {
	iter, err := s.leagues.FilterTeamTransfer(opts)
	if err != nil {
		return err
	}

	for iter.Next() {
		e := *(iter.Event)
		s.addEvent(e.Raw, "LeaguesTeamTransfer", e)
	}
	return nil
}

func (s *EventScanner) scanPlayerTransfer(opts *bind.FilterOpts) error {
	iter, err := s.leagues.FilterPlayerTransfer(opts)
	if err != nil {
		return err
	}

	for iter.Next() {
		e := *(iter.Event)
		s.addEvent(e.Raw, "LeaguesPlayerTransfer", e)
	}
	return nil
}

func (s *EventScanner) scanActionsSubmission(opts *bind.FilterOpts) error {
	iter, err := s.updates.FilterActionsSubmission(opts)
	if err != nil {
		return err
	}

	for iter.Next() {
		e := *(iter.Event)
		s.addEvent(e.Raw, "UpdatesActionsSubmission", e)
	}
	return nil
}
