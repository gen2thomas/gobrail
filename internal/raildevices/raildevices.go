package raildevices

import (
	"time"
)

// Timing is used for all kind of timing according to a rail device
type Timing struct {
	Starting time.Duration
	Stopping time.Duration
}

func limitTiming(timing Timing, maxTime time.Duration) Timing {
	if timing.Starting > maxTime {
		timing.Starting = maxTime
	}
	if timing.Stopping > maxTime {
		timing.Stopping = maxTime
	}
	return timing
}
