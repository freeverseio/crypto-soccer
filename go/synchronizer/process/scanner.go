package process

import (
	"context"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"

	"github.com/freeverseio/crypto-soccer/go/contracts"
)

type AbstractEvent struct {
	BlockNumber    uint64
	Timestamp      uint64
	TxIndexInBlock uint
	Name           string // for ease of debugging
	Value          interface{}
}

func NewAbstractEvent(blockNumber uint64, txIndex uint, name string, x interface{}) *AbstractEvent {
	return &AbstractEvent{blockNumber, 0, txIndex, name, x}
}

type EventScanner struct {
	contracts *contracts.Contracts
	Events    []*AbstractEvent
}

func NewEventScanner(contracts *contracts.Contracts) *EventScanner {
	return &EventScanner{contracts, []*AbstractEvent{}}
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
// assets.AssetsDivisionCreation
// assets.AssetsTeamTransfer
// assets.AssetsPlayerTransfer
// updates.UpdatesActionsSubmission

func (s *EventScanner) Process(opts *bind.FilterOpts) error {
	if err := s.scanAssetsInit(opts); err != nil {
		return err
	}
	if err := s.scanDivisionCreation(opts); err != nil {
		return err
	}
	if err := s.scanPlayerStateChange(opts); err != nil {
		return err
	}
	if err := s.scanTeamTransfer(opts); err != nil {
		return err
	}
	if err := s.scanActionsSubmission(opts); err != nil {
		return err
	}

	log.Debug("scanner got: ", len(s.Events), " Abstract Events")

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
	header, err := s.contracts.Client.HeaderByHash(context.Background(), rawEvent.BlockHash)
	if err != nil {
		log.Fatal(err)
	}
	abstractEvent := AbstractEvent{}
	abstractEvent.BlockNumber = rawEvent.BlockNumber
	abstractEvent.Name = name
	abstractEvent.Timestamp = header.Time
	abstractEvent.TxIndexInBlock = rawEvent.Index
	abstractEvent.Value = event
	s.Events = append(s.Events, &abstractEvent)
}

func (s *EventScanner) scanAssetsInit(opts *bind.FilterOpts) error {
	iter, err := s.contracts.Assets.FilterAssetsInit(opts)
	if err != nil {
		return err
	}
	for iter.Next() {
		e := *(iter.Event)
		log.Debugf("[scanner] scanAssetsInit")
		s.addEvent(e.Raw, "AssetsInit", e)
	}
	return nil
}

func (s *EventScanner) scanDivisionCreation(opts *bind.FilterOpts) error {
	iter, err := s.contracts.Assets.FilterDivisionCreation(opts)
	if err != nil {
		return err
	}
	for iter.Next() {
		e := *(iter.Event)
		log.Debugf("[scanner] scanDivisionCreation timezone %v", e.Timezone)
		s.addEvent(e.Raw, "AssetsDivisionCreation", e)
	}
	return nil
}

func (s *EventScanner) scanPlayerStateChange(opts *bind.FilterOpts) error {
	iter, err := s.contracts.Assets.FilterPlayerStateChange(opts)
	if err != nil {
		return err
	}
	for iter.Next() {
		e := *(iter.Event)
		log.Debugf("[scanner] scanPlayerStateChange playerId %v state %v", e.PlayerId, e.State)
		s.addEvent(e.Raw, "PlayerStateChange", e)
	}
	return nil
}

func (s *EventScanner) scanTeamTransfer(opts *bind.FilterOpts) error {
	iter, err := s.contracts.Assets.FilterTeamTransfer(opts)
	if err != nil {
		return err
	}
	for iter.Next() {
		e := *(iter.Event)
		log.Debugf("[scanner] scanTeamTransfer teamId %v to %v", e.TeamId, e.To.String())
		s.addEvent(e.Raw, "AssetsTeamTransfer", e)
	}
	return nil
}

//func (s *EventScanner) scanPlayerTransfer(opts *bind.FilterOpts) error {
//	for iter.Next() {
//		e := *(iter.Event)
//		log.Debugf("[scanner] scanPlayerTransfer playerId %v, toTeam %v", e.PlayerId, e.TeamIdTarget.String())
//		s.addEvent(e.Raw, "AssetsPlayerTransfer", e)
//	}
//	return nil
//}

func (s *EventScanner) scanPlayerFreeze(opts *bind.FilterOpts) error {
	iter, err := s.contracts.Market.FilterPlayerFreeze(opts)
	if err != nil {
		return err
	}
	for iter.Next() {
		e := *(iter.Event)
		log.Debugf("[scanner] scanPlayerFreeze")
		s.addEvent(e.Raw, "FilterPlayerFreeze", e)
	}
	return nil
}

func (s *EventScanner) scanActionsSubmission(opts *bind.FilterOpts) error {
	iter, err := s.contracts.Updates.FilterActionsSubmission(opts)
	if err != nil {
		return err
	}
	for iter.Next() {
		e := *(iter.Event)
		log.Debugf("[scanner] ScanActionSubmission")
		s.addEvent(e.Raw, "UpdatesActionsSubmission", e)
	}
	return nil
}
