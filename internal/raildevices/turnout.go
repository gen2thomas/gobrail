package raildevices

// A turnout or railroad switch is a rail device used for changing the direction of a train to a diverging route.
// The difference to a normal switch is the output must not be permant set to on, but only for a time period of 0.25-1s
// Second difference: A standard model turnout needs 2 physical outputs (left, right).

import (
	"fmt"
	"time"
)

const maxTime = time.Duration(time.Second)

// TurnoutDevice is describes a turnout
type TurnoutDevice struct {
	name             string
	nameBranch       string
	timing           Timing
	oldStateToBranch map[string]bool
	stateToBranch    bool
	boardsAPI        BoardsAPIer
	inputDevice      Inputer
	firstRun         bool
}

// NewTurnout creates an instance of a turnout
func NewTurnout(boardsAPI BoardsAPIer, boardID string, boardPinNr uint8, railDeviceName string, boardPinNrBranch uint8, timing Timing) (s *TurnoutDevice, err error) {
	if err = boardsAPI.MapBinaryPin(boardID, boardPinNr, railDeviceName); err != nil {
		return
	}
	railDeviceNameBranch := railDeviceName + " branch"
	if err = boardsAPI.MapBinaryPin(boardID, boardPinNrBranch, railDeviceNameBranch); err != nil {
		return
	}
	s = &TurnoutDevice{
		name:             railDeviceName,
		nameBranch:       railDeviceNameBranch,
		timing:           limitTiming(timing),
		oldStateToBranch: make(map[string]bool),
		boardsAPI:        boardsAPI,
	}
	return
}

// StateChanged states true when turnout status was changed since last visit
func (s *TurnoutDevice) StateChanged(visitor string) (hasChanged bool, err error) {
	oldStateToBranch, known := s.oldStateToBranch[visitor]
	if s.stateToBranch != oldStateToBranch || !known {
		s.oldStateToBranch[visitor] = s.stateToBranch
		hasChanged = true
	}
	return
}

// IsOn means the track switch is switched to this direction, that the train will run
// to the inner circle or diverging route from the main route
//
// =CHOO-CHOO>====== --> IsOn = false (train runs the main route or outer circle)
//              \\
//               ||  --> IsOn = true  (train runs the inner circle)
//              //
func (s *TurnoutDevice) IsOn() bool {
	return s.stateToBranch
}

// SwitchOn will try to switch on the turnout
func (s *TurnoutDevice) SwitchOn() (err error) {
	if err = s.boardsAPI.SetValue(s.nameBranch, 1); err != nil {
		return
	}
	time.Sleep(s.timing.Starting)
	if err = s.boardsAPI.SetValue(s.nameBranch, 0); err != nil {
		return
	}
	s.stateToBranch = true
	return
}

// SwitchOff will switch off the turnout
func (s *TurnoutDevice) SwitchOff() (err error) {
	if err = s.boardsAPI.SetValue(s.name, 1); err != nil {
		return
	}
	time.Sleep(s.timing.Stopping)
	if err = s.boardsAPI.SetValue(s.name, 0); err != nil {
		return
	}
	s.stateToBranch = false
	return
}

// Name gets the name of the turnout (rail device name)
func (s *TurnoutDevice) Name() string {
	return s.name
}

// Map is mapping an input for use in Run()
func (s *TurnoutDevice) Map(inputDevice Inputer) (err error) {
	if s.inputDevice != nil {
		return fmt.Errorf("turnout '%s' is already mapped to an input '%s'", s.name, s.inputDevice.Name())
	}
	if s.name == inputDevice.Name() {
		return fmt.Errorf("Circular mapping blocked for turnout '%s'", s.name)
	}
	s.inputDevice = inputDevice
	return nil
}

// Run is called in a loop and will make action dependant on the input device
func (s *TurnoutDevice) Run() (err error) {
	if s.inputDevice == nil {
		return fmt.Errorf("turnout '%s' can't run, please map to an input first", s.name)
	}
	var changed bool
	if changed, err = s.inputDevice.StateChanged(s.name); err != nil {
		return err
	}
	if !(changed || s.firstRun) {
		return
	}
	if s.inputDevice.IsOn() {
		s.SwitchOn()
	} else {
		s.SwitchOff()
	}
	return
}

// ReleaseInput is used to unmap
func (s *TurnoutDevice) ReleaseInput() {
	s.inputDevice = nil
}

func limitTiming(timing Timing) Timing {
	if timing.Starting > maxTime {
		timing.Starting = maxTime
	}
	if timing.Stopping > maxTime {
		timing.Stopping = maxTime
	}
	return timing
}
