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

// Timing is used for all kind of timing according to a rail device
type Timing struct {
	starting time.Duration
	stoping  time.Duration
}
