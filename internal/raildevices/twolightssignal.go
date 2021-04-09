package raildevices

// A two light signal is a rail device used for sign "pass" or "stop" with lamps (commonly in colours "green"/"red")
// Both lights can't be set at the same time.
// The output is static set in difference to a semaphore signal, which is more like a "turnout" to control.

import (
	"fmt"
	"time"
)

// TwoLightSignalDevice is describes a signal with two lights
type TwoLightSignalDevice struct {
	commonName     string
	name           string
	nameStop       string
	timing         Timing
	oldState       map[string]bool
	state          bool
	boardsAPI      BoardsAPIer
	inputDevice    Inputer
	inputInversion bool
	firstRun       bool
}

// NewTwoLightSignal creates an instance of a light signal with two lights
func NewTwoLightSignal(boardsAPI BoardsAPIer, boardID string, boardPinNr uint8, railDeviceName string, boardPinNrStop uint8, timing Timing) (s *TwoLightSignalDevice, err error) {
	if err = boardsAPI.MapBinaryPin(boardID, boardPinNr, railDeviceName); err != nil {
		return
	}
	railDeviceNameStop := railDeviceName + " stop"
	if err = boardsAPI.MapBinaryPin(boardID, boardPinNrStop, railDeviceNameStop); err != nil {
		return
	}
	s = &TwoLightSignalDevice{
		commonName: "two light signal",
		name:       railDeviceName,
		nameStop:   railDeviceNameStop,
		timing:     limitTiming(timing),
		oldState:   make(map[string]bool),
		boardsAPI:  boardsAPI,
	}
	return
}

// StateChanged states true when light signal status was changed since last visit
func (s *TwoLightSignalDevice) StateChanged(visitor string) (hasChanged bool, err error) {
	oldState, known := s.oldState[visitor]
	if s.state != oldState || !known {
		s.oldState[visitor] = s.state
		hasChanged = true
	}
	return
}

// IsOn means the "can pass" position (e.g. green colour)
func (s *TwoLightSignalDevice) IsOn() bool {
	return s.state
}

// SwitchOn will try to switch off the "stop" light (e.g. red colour)
// and immediatally switch on the "can pass" light (e.g. green colour)
func (s *TwoLightSignalDevice) SwitchOn() (err error) {
	time.Sleep(s.timing.Starting)
	if err = s.boardsAPI.SetValue(s.nameStop, 0); err != nil {
		return
	}
	if err = s.boardsAPI.SetValue(s.name, 1); err != nil {
		return
	}
	s.state = true
	return
}

// SwitchOff will try to switch off the "can pass" light (e.g. green colour)
// and immediatally switch on the "stop" light (e.g. red colour)
func (s *TwoLightSignalDevice) SwitchOff() (err error) {
	time.Sleep(s.timing.Stopping)
	if err = s.boardsAPI.SetValue(s.name, 0); err != nil {
		return
	}
	if err = s.boardsAPI.SetValue(s.nameStop, 1); err != nil {
		return
	}
	s.state = false
	return
}

// Name gets the name of the light signal (rail device name)
func (s *TwoLightSignalDevice) Name() string {
	return s.name
}

// Connect is connecting an input for use in Run()
func (s *TwoLightSignalDevice) Connect(inputDevice Inputer) (err error) {
	if s.inputDevice != nil {
		return fmt.Errorf("The %s '%s' is already mapped to an input '%s'", s.commonName, s.name, s.inputDevice.Name())
	}
	if s.name == inputDevice.Name() {
		return fmt.Errorf("Circular mapping blocked for %s '%s'", s.commonName, s.name)
	}
	s.inputDevice = inputDevice
	return nil
}

// ConnectInverse is connecting an input for use in Run(), but with inversed action
func (s *TwoLightSignalDevice) ConnectInverse(inputDevice Inputer) (err error) {
	s.Connect(inputDevice)
	s.inputInversion = true
	return nil
}

// Run is called in a loop and will make action dependant on the input device
func (s *TwoLightSignalDevice) Run() (err error) {
	if s.inputDevice == nil {
		return fmt.Errorf("The %s '%s' can't run, please map to an input first", s.commonName, s.name)
	}
	var changed bool
	if changed, err = s.inputDevice.StateChanged(s.name); err != nil {
		return err
	}
	if !(changed || s.firstRun) {
		return
	}
	if s.inputDevice.IsOn() != s.inputInversion {
		s.SwitchOn()
	} else {
		s.SwitchOff()
	}
	return
}

// ReleaseInput is used to unmap
func (s *TwoLightSignalDevice) ReleaseInput() {
	s.inputDevice = nil
}
