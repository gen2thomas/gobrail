package raildevices

// A common output is a rail device used for other output rail devices to minimize implementation and test effort

import (
	"fmt"
	"time"
)

// CommonOutputDevice describes a Common output device with one output
type CommonOutputDevice struct {
	label          string
	railDeviceName string
	timing         Timing
	oldState       map[string]bool
	state          bool
	defectiveState bool
	inputDevice    Inputer
	inputInversion bool
	firstRun       bool
}

// NewCommonOutput creates an instance of a rail device for usage with outputs
func NewCommonOutput(railDeviceName string, timing Timing, label string) (co *CommonOutputDevice) {
	co = &CommonOutputDevice{
		label:          label,
		railDeviceName: railDeviceName,
		timing:         timing,
		oldState:       make(map[string]bool),
		firstRun:       true,
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
		err = fmt.Errorf("The %s '%s' is defective, please repair before switch on", o.label, o.railDeviceName)
	}
	return
}

// MakeDefective causes the Common output device in an simulated defective state
func (o *CommonOutputDevice) MakeDefective(offFunc func() (err error)) (err error) {
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
		return fmt.Errorf("The %s '%s' can be only repaired when off", o.label, o.railDeviceName)
	}
	o.defectiveState = false
	return
}

// RailDeviceName gets the name of the Common output device
func (o *CommonOutputDevice) RailDeviceName() string {
	return o.railDeviceName
}

// Connect is connecting an input for use in Run()
func (o *CommonOutputDevice) Connect(inputDevice Inputer) (err error) {
	if o.inputDevice != nil {
		return fmt.Errorf("The %s '%s' is already connected to an input '%s'", o.label, o.railDeviceName, o.inputDevice.RailDeviceName())
	}
	if o.railDeviceName == inputDevice.RailDeviceName() {
		return fmt.Errorf("Circular mapping blocked for %s '%s'", o.label, o.railDeviceName)
	}
	o.inputDevice = inputDevice
	return nil
}

// ConnectInverse is connecting an input for use in Run(), but with inversed action
func (o *CommonOutputDevice) ConnectInverse(inputDevice Inputer) (err error) {
	o.Connect(inputDevice)
	o.inputInversion = true
	return nil
}

// Run is called in a loop and will make action dependant on the input device
func (o *CommonOutputDevice) Run(onFunc func() (err error), offFunc func() (err error)) (err error) {
	if o.inputDevice == nil {
		return fmt.Errorf("The %s '%s' can't run, please map to an input first", o.label, o.railDeviceName)
	}
	var changed bool
	if changed, err = o.inputDevice.StateChanged(o.railDeviceName); err != nil {
		return err
	}
	if !(changed || o.firstRun) {
		return
	}
	o.firstRun = false
	if o.inputDevice.IsOn() != o.inputInversion {
		err = onFunc()
	} else {
		err = offFunc()
	}
	return
}

// ReleaseInput is used to unmap
func (o *CommonOutputDevice) ReleaseInput() {
	o.inputDevice = nil
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
