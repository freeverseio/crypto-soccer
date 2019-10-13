package process

import (
	//"container/list"
	//"errors"
	"math/rand"
	//"sort"
	"strconv"
	"testing"
)

type testEvent1 struct {
	Minutes uint64
	Seconds uint
	X       int
	Y       int
	Name    string
}

type testEvent2 struct {
	Minutes uint64
	Seconds uint
	Value   string
	Name    string
}

func cast(x interface{}, test *testing.T) (*AbstractEvent, bool) {
	switch v := x.(type) {
	case *testEvent1:
		ev, _ := x.(*testEvent1)
		return &AbstractEvent{ev.Minutes, ev.Seconds, ev.Name, ev}, true
	case *testEvent2:
		ev, _ := x.(*testEvent2)
		return &AbstractEvent{ev.Minutes, ev.Seconds, ev.Name, ev}, true
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
	for i := 0; i < count; i++ {
		minutes := uint64(rand.Intn(60))
		seconds := uint(rand.Intn(60))
		if i&1 == 0 {
			l.Array[i] = &testEvent1{minutes, seconds, 1, 2, "event1"}
		} else {
			l.Array[i] = &testEvent2{minutes, seconds, strconv.Itoa(i), "event2"}
		}
	}
	return &l
}

func TestList(test *testing.T) {
	eventList := newEventList(10)
	genericEvents := []*AbstractEvent{}
	test.Log("========= creating events ==========")
	for i, v := range eventList.Array {
		test.Logf("adding event %v is %v", i, v)
		if e, ok := cast(v, test); ok {
			genericEvents = append(genericEvents, e)
		} else {
			test.Fatal("Unknown event type")
		}
	}

	time := func(p1, p2 *AbstractEvent) bool {
		if p1.BlockNumber == p2.BlockNumber {
			return p1.TxIndexInBlock < p2.TxIndexInBlock
		}
		return p1.BlockNumber < p2.BlockNumber

	}

	byFunction(time).Sort(genericEvents)

	test.Log("========= sorted events ===========")
	for i, v := range genericEvents {
		// casting back to original type
		switch v := v.Value.(type) {
		case *testEvent1:
			test.Logf("event %v is %v", i, *v)
		case *testEvent2:
			test.Logf("event %v is %v", i, *v)
		default:
			test.Fatalf("Could not cast type %v:", v)
		}
	}
}
