package raildevices

// A lamp is a rail device used for
// simple lamps, neon-light simulation, blinking lamps

import (
	"fmt"
	"time"
)

// LampDevice is describes a lamp
type LampDevice struct {
	name           string
	timing         Timing
	oldState       map[string]bool
	state          bool
	defectiveState bool
	boardsAPI      BoardsAPIer
	inputDevice    Inputer
	firstRun       bool
}

// NewLamp creates an instance of a lamp
func NewLamp(boardsAPI BoardsAPIer, boardID string, boardPinNr uint8, railDeviceName string, timing Timing) (ld *LampDevice, err error) {
	if err = boardsAPI.MapBinaryPin(boardID, boardPinNr, railDeviceName); err != nil {
		return
	}
	ld = &LampDevice{
		name:      railDeviceName,
		timing:    timing,
		oldState:  make(map[string]bool),
		boardsAPI: boardsAPI,
	}
	return
}

// StateChanged states true when lamp status was changed since last visit
func (l *LampDevice) StateChanged(visitor string) (hasChanged bool, err error) {
	oldState, known := l.oldState[visitor]
	if l.state != oldState || !known {
		l.oldState[visitor] = l.state
		hasChanged = true
	}
	return
}

// IsOn states true when lamp is on
func (l *LampDevice) IsOn() bool {
	return l.state
}

// IsDefective states true when lamp is defective
func (l *LampDevice) IsDefective() bool {
	return l.defectiveState
}

// SwitchOn will try to switch on the lamp
func (l *LampDevice) SwitchOn() (err error) {
	if l.IsDefective() {
		err = fmt.Errorf("Lamp '%s' is defective, please repair before switch on", l.name)
		return
	}
	time.Sleep(l.timing.Starting)
	if err = l.boardsAPI.SetValue(l.name, 1); err != nil {
		return
	}
	l.state = true
	return
}

// SwitchOff will switch off the lamp
func (l *LampDevice) SwitchOff() (err error) {
	time.Sleep(l.timing.Stopping)
	if err = l.boardsAPI.SetValue(l.name, 0); err != nil {
		return
	}
	l.state = false
	return
}

// MakeDefective causes the lamp in an simulated defective state
func (l *LampDevice) MakeDefective() (err error) {
	if err = l.SwitchOff(); err != nil {
		err = fmt.Errorf("Can't switch off before make defective, %w", err)
		return
	}
	l.defectiveState = true
	return
}

// Repair will fix the simulated defective state
func (l *LampDevice) Repair() (err error) {
	if l.IsOn() {
		return fmt.Errorf("Lamp '%s' can be only repaired when off", l.name)
	}
	l.defectiveState = false
	return
}

// Name gets the name of the lamp (rail device name)
func (l *LampDevice) Name() string {
	return l.name
}

// Map is mapping an input for use in Run()
func (l *LampDevice) Map(inputDevice Inputer) (err error) {
	if l.inputDevice != nil {
		return fmt.Errorf("Lamp '%s' is already mapped to an input '%s'", l.name, l.inputDevice.Name())
	}
	if l.name == inputDevice.Name() {
		return fmt.Errorf("Circular mapping blocked for Lamp '%s'", l.name)
	}
	l.inputDevice = inputDevice
	return nil
}

// Run is called in a loop and will make action dependant on the input device
func (l *LampDevice) Run() (err error) {
	if l.inputDevice == nil {
		return fmt.Errorf("Lamp '%s' can't run, please map to an input first", l.name)
	}
	var changed bool
	if changed, err = l.inputDevice.StateChanged(l.name); err != nil {
		return err
	}
	if !(changed || l.firstRun) {
		return
	}
	if l.inputDevice.IsOn() {
		l.SwitchOn()
	} else {
		l.SwitchOff()
	}
	return
}

// ReleaseInput is used to unmap
func (l *LampDevice) ReleaseInput() {
	l.inputDevice = nil
}
