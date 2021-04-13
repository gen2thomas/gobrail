package raildevices

// A common output is a rail device used for other output rail devices to minimize implementation and test effort

import (
	"fmt"
	"time"
)

// CommonOutputDevice describes a Common output device with one output
type CommonOutputDevice struct {
	railDeviceName string
	timing         Timing
	oldState       map[string]bool
	state          bool
	defectiveState bool
}

// NewCommonOutput creates an instance of a rail device for usage with outputs
func NewCommonOutput(railDeviceName string, timing Timing) (co *CommonOutputDevice) {
	co = &CommonOutputDevice{
		railDeviceName: railDeviceName,
		timing:         timing,
		oldState:       make(map[string]bool),
	}
	return
}

// StateChanged states true when Common output device status was changed since last visit
func (o *CommonOutputDevice) StateChanged(visitor string) (hasChanged bool, err error) {
	oldState, known := o.oldState[visitor]
	if o.state != oldState || !known {
		o.oldState[visitor] = o.state
		hasChanged = true
	}
	return
}

// IsOn states true when Common output device is on
func (o *CommonOutputDevice) IsOn() bool {
	return o.state
}

// IsDefective states true when Common output device is defective
func (o *CommonOutputDevice) IsDefective() (err error) {
	if o.defectiveState {
		err = fmt.Errorf("The '%s' is defective, please repair before switch on", o.railDeviceName)
	}
	return
}

// MakeDefectiveCommon causes the Common output device in an simulated defective state
func (o *CommonOutputDevice) MakeDefectiveCommon(offFunc func() (err error)) (err error) {
	if err = offFunc(); err != nil {
		err = fmt.Errorf("Can't switch off before make defective, %w", err)
		return
	}
	o.defectiveState = true
	return
}

// Repair will fix the simulated defective state
func (o *CommonOutputDevice) Repair() (err error) {
	if o.IsOn() {
		return fmt.Errorf("The '%s' can be only repaired when off", o.railDeviceName)
	}
	o.defectiveState = false
	return
}

// RailDeviceName gets the name of the Common output device
func (o *CommonOutputDevice) RailDeviceName() string {
	return o.railDeviceName
}

// TimingForStart execute sleep with starttime
func (o *CommonOutputDevice) TimingForStart() {
	time.Sleep(o.timing.Starting)
}

// TimingForStop execute sleep with stoptime
func (o *CommonOutputDevice) TimingForStop() {
	time.Sleep(o.timing.Stopping)
}

// SetState sets the new state
func (o *CommonOutputDevice) SetState(newState bool) {
	o.state = newState
}
