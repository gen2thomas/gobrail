package raildevices

// A lamp is a rail device used for
// simple lamps, neon-light simulation, blinking lamps

import (
	"github.com/gen2thomas/gobrail/internal/boardpin"
)

// LampDevice is describes a lamp
type LampDevice struct {
	*CommonOutputDevice
	output *boardpin.Output
}

// NewLamp creates an instance of a lamp
func NewLamp(co *CommonOutputDevice, output *boardpin.Output) (ld *LampDevice) {
	ld = &LampDevice{
		CommonOutputDevice: co,
		output:             output,
	}
	return
}

// SwitchOn will try to switch on the StdLamp
func (l *LampDevice) SwitchOn() (err error) {
	if err = l.IsDefective(); err != nil {
		return
	}
	l.TimingForStart()
	if err = l.output.WriteValue(1); err != nil {
		return
	}
	l.SetState(true)
	return
}

// SwitchOff will switch off the StdLamp
func (l *LampDevice) SwitchOff() (err error) {
	l.TimingForStop()
	if err = l.output.WriteValue(0); err != nil {
		return
	}
	l.SetState(false)
	return
}

// MakeDefective causes the StdLamp in an simulated defective state
func (l *LampDevice) MakeDefective() (err error) {
	return l.MakeDefectiveCommon(l.SwitchOff)
}
