package process

import (
	//"container/list"
	//"errors"
	"math/rand"
	"sort"
	"strconv"
	"testing"
)

type mytime struct {
	Minutes int
	Seconds int
}

type event1 struct {
	Time mytime
	X    int
	Y    int
	Name string
}

type event2 struct {
	Time  mytime
	Value string
	Name  string
}
type genericEvent struct {
	Time  mytime
	Value interface{}
}

func newEvent1(x int, y int, t mytime) *event1 {
	return &event1{t, x, y, "event1"}
}

func newEvent2(value string, t mytime) *event2 {
	return &event2{t, value, "event2"}
}

func newGenericEvent(t mytime, x interface{}) *genericEvent {
	return &genericEvent{t, x}
}

func cast(x interface{}, test *testing.T) (*genericEvent, bool) {
	switch v := x.(type) {
	case *event1:
		ev, _ := x.(*event1)
		return newGenericEvent(ev.Time, ev), true
	case *event2:
		ev, _ := x.(*event2)
		return newGenericEvent(ev.Time, ev), true
	default:
		test.Logf("Could not cast type %v:", v)
		return nil, false
	}
}

type eventList struct {
	Array []interface{}
}

func newEventList(count int) *eventList {
	l := eventList{make([]interface{}, count)}
	t := mytime{0, 0}
	for i := 0; i < count; i++ {
		t.Minutes = rand.Intn(60)
		t.Seconds = rand.Intn(60)
		if i&1 == 0 {
			l.Array[i] = newEvent1(i, i*2, t)
		} else {
			l.Array[i] = newEvent2(strconv.Itoa(i), t)
		}
	}
	return &l
}

type by func(p1, p2 *genericEvent) bool

func (b by) Sort(events []*genericEvent) {
	ps := &eventSorter{
		events: events,
		by:     b,
	}
	sort.Sort(ps)
}

type eventSorter struct {
	events []*genericEvent
	by     func(p1, p2 *genericEvent) bool // Closure used in the Less method.
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
	eventList := newEventList(10)
	genericEvents := []*genericEvent{}
	test.Log("========= creating events ==========")
	for i, v := range eventList.Array {
		test.Logf("adding event %v is %v", i, v)
		if e, ok := cast(v, test); ok {
			genericEvents = append(genericEvents, e)
		} else {
			test.Fatal("Unknown event type")
		}
	}

	mytime := func(p1, p2 *genericEvent) bool {
		if p1.Time.Minutes == p2.Time.Minutes {
			return p1.Time.Seconds < p2.Time.Seconds
		}
		return p1.Time.Minutes < p2.Time.Minutes

	}

	by(mytime).Sort(genericEvents)

	test.Log("========= sorted events ===========")
	for i, v := range genericEvents {
		// casting back to original type
		switch v := v.Value.(type) {
		case *event1:
			test.Logf("event %v is %v", i, *v)
		case *event2:
			test.Logf("event %v is %v", i, *v)
		default:
			test.Fatalf("Could not cast type %v:", v)
		}
	}
}
