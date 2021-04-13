package raildevices

// A two light signal is a rail device used for sign "pass" or "stop" with lamps (commonly in colours "green"/"red")
// Both lights can't be set at the same time.
// The output pin is static set in difference to a semaphore signal, which is more like a "turnout" to control.

import (
	"github.com/gen2thomas/gobrail/internal/boardpin"
)

// TwoLightsSignalDevice is describes a signal with two lights
type TwoLightsSignalDevice struct {
	outputPass *boardpin.Output
	outputStop *boardpin.Output
	*CommonOutputDevice
}

// NewTwoLightsSignal creates an instance of a light signal with two lights
func NewTwoLightsSignal(co *CommonOutputDevice, outputPass *boardpin.Output, outputStop *boardpin.Output) (s *TwoLightsSignalDevice) {
	s = &TwoLightsSignalDevice{
		outputPass:         outputPass,
		outputStop:         outputStop,
		CommonOutputDevice: co,
	}
	return
}

// SwitchOn will try to switch off the "stop" light (e.g. red colour)
// and immediatally switch on the "can pass" light (e.g. green colour)
func (s *TwoLightsSignalDevice) SwitchOn() (err error) {
	s.TimingForStart()
	if err = s.outputStop.WriteValue(0); err != nil {
		return
	}
	if err = s.outputPass.WriteValue(1); err != nil {
		return
	}
	s.SetState(true)
	return
}

// SwitchOff will try to switch off the "can pass" light (e.g. green colour)
// and immediatally switch on the "stop" light (e.g. red colour)
func (s *TwoLightsSignalDevice) SwitchOff() (err error) {
	s.TimingForStop()
	if err = s.outputPass.WriteValue(0); err != nil {
		return
	}
	if err = s.outputStop.WriteValue(1); err != nil {
		return
	}
	s.SetState(false)
	return
}
