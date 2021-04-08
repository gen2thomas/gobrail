package raildevices

// A ToggleButton is a rail device used for an input by a button
// the output will change on each press of button

import (
	"fmt"
)

// ToggleButtonDevice describes a ToggleButton
type ToggleButtonDevice struct {
	name        string
	oldState    bool
	toggleState bool
	boardsAPI   BoardsAPIer
}

// NewToggleButton creates an instance of a ToggleButton
func NewToggleButton(boardsAPI BoardsAPIer, boardID string, boardPinNr uint8, railDeviceName string) (ld *ToggleButtonDevice, err error) {
	if err = boardsAPI.MapBinaryPin(boardID, boardPinNr, railDeviceName); err != nil {
		return
	}
	ld = &ToggleButtonDevice{
		name:      railDeviceName,
		boardsAPI: boardsAPI,
	}
	return
}

// StateChanged states true when ToggleButton status was changed
func (b *ToggleButtonDevice) StateChanged() (hasChanged bool, err error) {
	var value uint8
	if value, err = b.boardsAPI.GetValue(b.name); err != nil {
		err = fmt.Errorf("Can't read value from '%s', %w", b.name, err)
		return
	}
	newState := value > 0
	if !b.oldState && newState {
		b.toggleState = !b.toggleState
		hasChanged = true
	}
	b.oldState = newState
	return
}

// IsOn states true when toggle state is on
func (b *ToggleButtonDevice) IsOn() bool {
	return b.toggleState
}

// Name gets the name of the ToggleButton (rail device name)
func (b *ToggleButtonDevice) Name() string {
	return b.name
}
