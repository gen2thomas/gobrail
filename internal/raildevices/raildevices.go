package raildevices

import (
	"time"
)

// BoardsAPIer is an interface for interact with a boards API
type BoardsAPIer interface {
	MapBinaryPin(boardID string, boardPinNr uint8, railDeviceName string) (err error)
	MapAnalogPin(boardID string, boardPinNr uint8, railDeviceName string) (err error)
	MapMemoryPin(boardID string, boardPinNrOrNegative int, railDeviceName string) (err error)
	GetValue(railDeviceName string) (value uint8, err error)
	SetValue(railDeviceName string, value uint8) (err error)
}

// Inputer is an interface for input devices to map in output devices. When an output device
// have this functions it can be used as input for an successive device.
type Inputer interface {
	Name() string
	StateChanged() (hasChanged bool, err error)
	IsOn() bool
}

// Outputer is an interface for output devices
type Outputer interface {
	Name() string
	// Map is used to map an input for action (IsOn --> e.g. SwitchOn)
	Map(input Inputer) (err error)
	// Run must be called in a loop
	Run() (err error)
	// ReleaseInput is used to unmap the input
	ReleaseInput()
}

// Timing is used for all kind of timing according to a rail device
type Timing struct {
	starting time.Duration
	stoping  time.Duration
}
