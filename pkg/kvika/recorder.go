package kvika

import (
	"sort"
	"time"
)

type Event struct {
	At      float64
	Name    string
	Payload interface{}
}

type Events []Event

type Recorder struct {
	begin  time.Time
	events []Event
}

func newRecorder() *Recorder {
	return &Recorder{
		begin:  time.Now(),
		events: make([]Event, 0),
	}
}

func (r *Recorder) Start() {
	r.begin = time.Now()
}

func (r *Recorder) Record(name string, payload interface{}) *Recorder {
	r.events = append(r.events, Event{
		At:      float64(time.Now().Sub(r.begin).Nanoseconds()) / 1000000.0,
		Name:    name,
		Payload: payload,
	})
	return r
}

func (r *Recorder) recordRaw(at float64, name string, payload interface{}) *Recorder {
	r.events = append(r.events, Event{
		At:      at,
		Name:    name,
		Payload: payload,
	})
	return r
}

func (r *Recorder) sortedEvents() Events {
	sort.Slice(r.events, func(i, j int) bool {
		return r.events[i].At < r.events[j].At
	})
	return Events(r.events[:])
}
