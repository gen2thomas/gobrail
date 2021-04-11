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
	cmnOutDev  *CommonOutputDevice
}

// NewTwoLightsSignal creates an instance of a light signal with two lights
func NewTwoLightsSignal(co *CommonOutputDevice, outputPass *boardpin.Output, outputStop *boardpin.Output) (s *TwoLightsSignalDevice) {
	s = &TwoLightsSignalDevice{
		outputPass: outputPass,
		outputStop: outputStop,
		cmnOutDev:  co,
	}
	return
}

// SwitchOn will try to switch off the "stop" light (e.g. red colour)
// and immediatally switch on the "can pass" light (e.g. green colour)
func (s *TwoLightsSignalDevice) SwitchOn() (err error) {
	s.cmnOutDev.TimingForStart()
	if err = s.outputStop.WriteValue(0); err != nil {
		return
	}
	if err = s.outputPass.WriteValue(1); err != nil {
		return
	}
	s.cmnOutDev.SetState(true)
	return
}

// SwitchOff will try to switch off the "can pass" light (e.g. green colour)
// and immediatally switch on the "stop" light (e.g. red colour)
func (s *TwoLightsSignalDevice) SwitchOff() (err error) {
	s.cmnOutDev.TimingForStop()
	if err = s.outputPass.WriteValue(0); err != nil {
		return
	}
	if err = s.outputStop.WriteValue(1); err != nil {
		return
	}
	s.cmnOutDev.SetState(false)
	return
}

// StateChanged states true when light signal status was changed since last visit
func (s *TwoLightsSignalDevice) StateChanged(visitor string) (hasChanged bool, err error) {
	return s.cmnOutDev.StateChanged(visitor)
}

// IsOn means the "can pass" position (e.g. green colour)
func (s *TwoLightsSignalDevice) IsOn() bool {
	return s.cmnOutDev.IsOn()
}

// RailDeviceName gets the name of the light signal common output
func (s *TwoLightsSignalDevice) RailDeviceName() string {
	return s.cmnOutDev.RailDeviceName()
}

// Connect is connecting an input for use in Run()
func (s *TwoLightsSignalDevice) Connect(inputDevice Inputer) (err error) {
	return s.cmnOutDev.Connect(inputDevice)
}

// ConnectInverse is connecting an input for use in Run(), but with inversed action
func (s *TwoLightsSignalDevice) ConnectInverse(inputDevice Inputer) (err error) {
	return s.cmnOutDev.ConnectInverse(inputDevice)
}

// Run is called in a loop and will make action dependant on the input device
func (s *TwoLightsSignalDevice) Run() (err error) {
	return s.cmnOutDev.Run(s.SwitchOn, s.SwitchOff)
}

// ReleaseInput is used to unmap
func (s *TwoLightsSignalDevice) ReleaseInput() {
	s.cmnOutDev.ReleaseInput()
}
