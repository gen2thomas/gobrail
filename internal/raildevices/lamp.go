package raildevices

// A lamp is a rail device used for
// simple lamps, neon-light simulation, blinking lamps

import (
	"github.com/gen2thomas/gobrail/internal/boardpin"
)

// LampDevice is describes a lamp
type LampDevice struct {
	cmnOutDev *CommonOutputDevice
	output    *boardpin.Output
}

// NewLamp creates an instance of a lamp
func NewLamp(co *CommonOutputDevice, output *boardpin.Output) (ld *LampDevice) {
	ld = &LampDevice{
		cmnOutDev: co,
		output:    output,
	}
	return
}

// StateChanged states true when StdLamp status was changed since last visit
func (l *LampDevice) StateChanged(visitor string) (hasChanged bool, err error) {
	return l.cmnOutDev.StateChanged(visitor)
}

// IsOn states true when StdLamp is on
func (l *LampDevice) IsOn() bool {
	return l.cmnOutDev.IsOn()
}

// IsDefective states true when StdLamp is defective
func (l *LampDevice) IsDefective() (err error) {
	return l.cmnOutDev.IsDefective()
}

// SwitchOn will try to switch on the StdLamp
func (l *LampDevice) SwitchOn() (err error) {
	if err = l.cmnOutDev.IsDefective(); err != nil {
		return
	}
	l.cmnOutDev.TimingForStart()
	if err = l.output.WriteValue(1); err != nil {
		return
	}
	l.cmnOutDev.SetState(true)
	return
}

// SwitchOff will switch off the StdLamp
func (l *LampDevice) SwitchOff() (err error) {
	l.cmnOutDev.TimingForStop()
	if err = l.output.WriteValue(0); err != nil {
		return
	}
	l.cmnOutDev.SetState(false)
	return
}

// MakeDefective causes the StdLamp in an simulated defective state
func (l *LampDevice) MakeDefective() (err error) {
	return l.cmnOutDev.MakeDefective(l.SwitchOff)
}

// Repair will fix the simulated defective state
func (l *LampDevice) Repair() (err error) {
	return l.cmnOutDev.Repair()
}

// RailDeviceName gets the name of the lamp output
func (l *LampDevice) RailDeviceName() string {
	return l.cmnOutDev.RailDeviceName()
}

// Connect is connecting an input for use in Run()
func (l *LampDevice) Connect(inputDevice Inputer) (err error) {
	return l.cmnOutDev.Connect(inputDevice)
}

// ConnectInverse is connecting an input for use in Run(), but with inversed action
func (l *LampDevice) ConnectInverse(inputDevice Inputer) (err error) {
	return l.cmnOutDev.ConnectInverse(inputDevice)
}

// Run is called in a loop and will make action dependant on the input device
func (l *LampDevice) Run() (err error) {
	return l.cmnOutDev.Run(l.SwitchOn, l.SwitchOff)
}

// ReleaseInput is used to unmap
func (l *LampDevice) ReleaseInput() {
	l.cmnOutDev.ReleaseInput()
}
