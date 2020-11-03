package kvika

import (
	"testing"
)

func TestSortedEvents(t *testing.T) {
	r := newRecorder()
	r.recordRaw(1, "1", 1)
	r.recordRaw(3, "3", 3)
	r.recordRaw(2, "2", 2)
	events := r.sortedEvents()
	for i, ev := range events {
		if (i + 1) != ev.Payload.(int) {
			t.Fatalf("Event is not sorted. Expected=%v, Actual=%v", i+1, ev.Payload)
		}
	}
}
