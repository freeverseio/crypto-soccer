package process

import (
	//"container/list"
	//"errors"
	"math/rand"
	"sort"
	"testing"
)

type Time struct {
	Minutes int
	Seconds int
}

func (t *Time) Add() {

	t.Seconds += 1
	if t.Seconds >= 5 {
		t.Seconds = 0
		t.Minutes += 1
	}
}

type Event1 struct {
	X    int
	Y    int
	Name string
	Time Time
}

func NewEvent1(x int, y int, t Time) *Event1 {
	return &Event1{x, y, "event1", t}
}

func NewEvent2(value string, t Time) *Event2 {
	return &Event2{value, "Event2", t}
}

type Event2 struct {
	Value string
	Name  string
	Time  Time
}
type GenericEvent struct {
	Time  Time
	Value interface{}
}

func NewGenericEvent(t Time, x interface{}) *GenericEvent {
	return &GenericEvent{t, x}
}

func Cast(x interface{}, test *testing.T) (*GenericEvent, bool) {
	switch v := x.(type) {
	case *Event1:
		ev, _ := x.(*Event1)
		return NewGenericEvent(ev.Time, ev), true
	case *Event2:
		ev, _ := x.(*Event2)
		return NewGenericEvent(ev.Time, ev), true
	default:
		test.Logf("Could not cast type %v:", v)
		return nil, false
	}
}

type EventList struct {
	Array []interface{}
}

func NewEventList(count int) *EventList {
	l := EventList{make([]interface{}, count)}
	t := Time{0, 0}
	for i := 0; i < count; i++ {
		t.Minutes = rand.Intn(60)
		t.Seconds = rand.Intn(60)
		if i&1 == 0 {
			l.Array[i] = NewEvent1(i, i*2, t)
		} else {
			l.Array[i] = NewEvent2(string(i), t)
		}
	}
	return &l
}

type By func(p1, p2 *GenericEvent) bool

func (by By) Sort(events []*GenericEvent) {
	ps := &eventSorter{
		events: events,
		by:     by,
	}
	sort.Sort(ps)
}

type eventSorter struct {
	events []*GenericEvent
	by     func(p1, p2 *GenericEvent) bool // Closure used in the Less method.
}

func (s *eventSorter) Len() int {
	return len(s.events)
}

func (s *eventSorter) Swap(i, j int) {
	s.events[i], s.events[j] = s.events[j], s.events[i]
}

func (s *eventSorter) Less(i, j int) bool {
	return s.by(s.events[i], s.events[j])
}

func TestList(test *testing.T) {
	eventList := NewEventList(10)
	genericEvents := []*GenericEvent{}
	test.Log("========= creating events ==========")
	for i, v := range eventList.Array {
		test.Logf("event %v is %v", i, v)
		if e, ok := Cast(v, test); ok {
			genericEvents = append(genericEvents, e)
		} else {
			test.Fatal("Unknown event type")
		}
	}

	time := func(p1, p2 *GenericEvent) bool {
		if p1.Time.Minutes == p2.Time.Minutes {
			return p1.Time.Seconds < p2.Time.Seconds
		}
		return p1.Time.Minutes < p2.Time.Minutes

	}

	By(time).Sort(genericEvents)

	test.Log("========= sorted events ===========")
	for i, v := range genericEvents {
		test.Logf("event %v is %v", i, v)
	}
}
