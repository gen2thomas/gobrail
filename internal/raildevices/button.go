package raildevices

// A Button is a rail device used for an input by button

import (
	"fmt"
)

// ButtonDevice describes a Button
type ButtonDevice struct {
	railDeviceName string
	state          bool
	oldState       map[string]bool
	boardsAPI      BoardsAPIer
}

// NewButton creates an instance of a Button
func NewButton(boardsAPI BoardsAPIer, boardID string, boardPinNr uint8, railDeviceName string) (b *ButtonDevice, err error) {
	if err = boardsAPI.MapBinaryPin(boardID, boardPinNr, railDeviceName); err != nil {
		return
	}
	b = &ButtonDevice{
		railDeviceName: railDeviceName,
		oldState:       make(map[string]bool),
		boardsAPI:      boardsAPI,
	}
	return
}

// StateChanged states true when Button status was changed
func (b *ButtonDevice) StateChanged(visitor string) (hasChanged bool, err error) {
	var value uint8
	if value, err = b.boardsAPI.GetValue(b.railDeviceName); err != nil {
		err = fmt.Errorf("Can't read value from '%s', %w", b.railDeviceName, err)
		return
	}
	b.state = value > 0
	oldState, known := b.oldState[visitor]
	if b.state != oldState || !known {
		b.oldState[visitor] = b.state
		hasChanged = true
	}
	return
}

// IsOn gets the state of the button
func (b *ButtonDevice) IsOn() bool {
	return b.state
}

// RailDeviceName gets the name of the button input
func (b *ButtonDevice) RailDeviceName() string {
	return b.railDeviceName
}
