package raildevices

// A lamp is a rail device used for
// simple lamps, neon-light simulation, blinking lamps

import (
	"fmt"
	"time"
)

// LampDevice is describes a lamp
type LampDevice struct {
	name          string
	stateName     string
	defectiveName string
	timing        Timing
	boardsAPI     BoardsAPIer
	inputDevice   Inputer
}

// NewLamp creates an instance of a lamp
func NewLamp(boardsAPI BoardsAPIer, boardID string, boardPinNr uint8, railDeviceName string, timing Timing) (ld *LampDevice, err error) {
	stateName := railDeviceName + " state"
	defectiveName := railDeviceName + " defective"
	if err = boardsAPI.MapBinaryPin(boardID, boardPinNr, railDeviceName); err != nil {
		return
	}
	if err = boardsAPI.MapMemoryPin(boardID, -1, stateName); err != nil {
		return
	}
	if err = boardsAPI.MapMemoryPin(boardID, -1, defectiveName); err != nil {
		return
	}
	ld = &LampDevice{
		name:          railDeviceName,
		stateName:     stateName,
		defectiveName: defectiveName,
		boardsAPI:     boardsAPI,
	}
	if err = ld.SwitchOff(); err != nil {
		return
	}
	if err = ld.Repair(); err != nil {
		return
	}
	return
}

// IsOn states true when lamp is on
func (l *LampDevice) IsOn() (isOn bool, err error) {
	var value uint8
	if value, err = l.boardsAPI.GetValue(l.stateName); err != nil {
		err = fmt.Errorf("Can't read value from '%s', %w", l.stateName, err)
		return
	}
	return value > 0, nil
}

// IsDefective states true when lamp is defective
func (l *LampDevice) IsDefective() (isDefect bool, err error) {
	var value uint8
	if value, err = l.boardsAPI.GetValue(l.defectiveName); err != nil {
		err = fmt.Errorf("Can't read value from '%s', %w", l.defectiveName, err)
		return
	}
	isDefect = value > 0
	return
}

// SwitchOn will try to switch on the lamp
func (l *LampDevice) SwitchOn() (err error) {
	var isDefect bool
	if isDefect, err = l.IsDefective(); err != nil {
		err = fmt.Errorf("Can't detect defective state before switch on, %w", err)
		return
	}
	if isDefect {
		err = fmt.Errorf("Lamp '%s' is defective, please repair before switch on", l.name)
		return
	}
	time.Sleep(l.timing.starting)
	if err = l.boardsAPI.SetValue(l.name, 1); err != nil {
		return
	}
	return l.boardsAPI.SetValue(l.stateName, 1)
}

// SwitchOff will switch off the lamp
func (l *LampDevice) SwitchOff() (err error) {
	time.Sleep(l.timing.stoping)
	if err = l.boardsAPI.SetValue(l.name, 0); err != nil {
		return
	}
	return l.boardsAPI.SetValue(l.stateName, 0)
}

// MakeDefective causes the lamp in an simulated defective state
func (l *LampDevice) MakeDefective() (err error) {
	if err = l.SwitchOff(); err != nil {
		err = fmt.Errorf("Can't switch off before make defective, %w", err)
		return
	}
	return l.boardsAPI.SetValue(l.defectiveName, 1)
}

// Repair will fix the simulated defective state
func (l *LampDevice) Repair() (err error) {
	var isOn bool
	if isOn, err = l.IsOn(); err != nil {
		return err
	}
	if isOn {
		return fmt.Errorf("Lamp '%s' can be only repaired when off", l.name)
	}
	return l.boardsAPI.SetValue(l.defectiveName, 0)
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
	// synchronize first
	if l.inputDevice.IsOn() {
		l.SwitchOn()
	} else {
		l.SwitchOff()
	}
	return nil
}

// Run is called in a loop and will make action dependant on the input device
func (l *LampDevice) Run() (err error) {
	if l.inputDevice == nil {
		return fmt.Errorf("Lamp '%s' can't run, please map to an input first", l.name)
	}
	var changed bool
	if changed, err = l.inputDevice.StateChanged(); err != nil {
		return err
	}
	if !changed {
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
