package raildevices

import (
	"time"
)

// Inputer is an interface for input devices to map in output devices. When an output device
// have this functions it can be used as input for an successive device.
type Inputer interface {
	RailDeviceName() string
	StateChanged(visitor string) (hasChanged bool, err error)
	IsOn() bool
}

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
