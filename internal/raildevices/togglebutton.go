package raildevices

// A ToggleButton is a rail device used for an input by a button
// the output will change on each press of button

import (
	"fmt"
)

// ToggleButtonDevice describes a ToggleButton
type ToggleButtonDevice struct {
	railDeviceName string
	oldState       bool
	toggleState    bool
	oldToggleState map[string]bool
	boardsAPI      BoardsAPIer
}

// NewToggleButton creates an instance of a ToggleButton
func NewToggleButton(boardsAPI BoardsAPIer, boardID string, boardPinNr uint8, railDeviceName string) (ld *ToggleButtonDevice, err error) {
	if err = boardsAPI.MapBinaryPin(boardID, boardPinNr, railDeviceName); err != nil {
		return
	}
	ld = &ToggleButtonDevice{
		railDeviceName: railDeviceName,
		oldToggleState: make(map[string]bool),
		boardsAPI:      boardsAPI,
	}
	return
}

// StateChanged states true when ToggleButton status was changed
func (b *ToggleButtonDevice) StateChanged(visitor string) (hasChanged bool, err error) {
	var value uint8
	if value, err = b.boardsAPI.GetValue(b.railDeviceName); err != nil {
		err = fmt.Errorf("Can't read value from '%s', %w", b.railDeviceName, err)
		return
	}
	// toggle button
	newState := value > 0
	if !b.oldState && newState {
		b.toggleState = !b.toggleState
	}
	b.oldState = newState
	// visitor
	oldToggleState, known := b.oldToggleState[visitor]
	if b.toggleState != oldToggleState || !known {
		hasChanged = true
		b.oldToggleState[visitor] = b.toggleState
	}
	return
}

// IsOn states true when toggle state is on
func (b *ToggleButtonDevice) IsOn() bool {
	return b.toggleState
}

// RailDeviceName gets the name of the toggle button input
func (b *ToggleButtonDevice) RailDeviceName() string {
	return b.railDeviceName
}
