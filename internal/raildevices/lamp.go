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
	boardsAPI     boardsAPIer
	timing        Timing
}

// NewLamp creates an instance of a lamp
func NewLamp(boardsAPI boardsApier, boardID string, boardPinNr uint8, railDeviceName string, timing Timing) *LampDevice {
	stateName := railDeviceName + " state"
	defectiveName := railDeviceName + " defective"

	boardsAPI.MapBinaryPin(boardID, boardPinNr, railDeviceName)
	boardsAPI.MapMemoryPin(boardID, -1, stateName)
	boardsAPI.MapMemoryPin(boardID, -1, defectiveName)
	ld := &LampDevice{
		name:          railDeviceName,
		stateName:     stateName,
		defectiveName: defectiveName,
		boardsAPI:     boardsAPI,
	}
	ld.SwitchOff()
	ld.Repair()
	return ld
}

// IsOn states true when lamp is on
func (l *LampDevice) IsOn() bool {
	value, err := l.boardsAPI.GetValue(l.stateName)
	if err != nil {
		fmt.Printf("Can't read value from '%s', %s\n", l.stateName, err)
		return false
	}
	return value > 0
}

// IsOff states true when lamp is off
func (l *LampDevice) IsOff() bool {
	return !l.IsOn()
}

// IsDefective states true when lamp is defective
func (l *LampDevice) IsDefective() bool {
	value, err := l.boardsAPI.GetValue(l.defectiveName)
	if err != nil {
		fmt.Printf("Can't read value from '%s', %s\n", l.defectiveName, err)
		return false
	}
	defective := value > 0
	if !defective {
		fmt.Printf("Lamp '%s' is working\n", l.name)
	}
	return defective
}

// SwitchOn will try to switch on the lamp
func (l *LampDevice) SwitchOn() {
	if l.IsDefective() {
		fmt.Printf("Lamp '%s' is defective, please repair before switch on\n", l.name)
		return
	}
	time.Sleep(l.timing.starting)
	l.boardsAPI.SetValue(l.name, 1)
	l.boardsAPI.SetValue(l.stateName, 1)
}

// SwitchOff will switch off the lamp
func (l *LampDevice) SwitchOff() {
	time.Sleep(l.timing.stoping)
	l.boardsAPI.SetValue(l.name, 0)
	l.boardsAPI.SetValue(l.stateName, 0)
}

// MakeDefective causes the lamp in an simulated defective state
func (l *LampDevice) MakeDefective() {
	l.SwitchOff()
	l.boardsAPI.SetValue(l.defectiveName, 1)
	fmt.Printf("Lamp '%s' is now defective, please repair\n", l.name)
}

// Repair will fix the simulated defective state
func (l *LampDevice) Repair() {
	if l.IsOn() {
		fmt.Printf("Lamp '%s' can be only repaired when off\n", l.name)
		return
	}
	fmt.Printf("Lamp '%s' is working again\n", l.name)
	l.boardsAPI.SetValue(l.defectiveName, 0)
}
