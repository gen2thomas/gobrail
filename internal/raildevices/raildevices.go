package raildevices

import (
	"time"
)

// Timing is used for all kind of timing according to a rail device
type Timing struct {
	Starting time.Duration
	Stopping time.Duration
}

// Limit is used to shrink start and stop time to the given maximum
func (t *Timing) Limit(maxTime time.Duration) {
	if t.Starting > maxTime {
		t.Starting = maxTime
	}
	if t.Stopping > maxTime {
		t.Stopping = maxTime
	}
}
