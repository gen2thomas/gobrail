package raildevices

// A Button is a rail device used for an input by button

import (
	"fmt"
)

// ButtonDevice describes a Button
type ButtonDevice struct {
	name      string
	oldState  bool
	boardsAPI BoardsAPIer
}

// NewButton creates an instance of a Button
func NewButton(boardsAPI BoardsAPIer, boardID string, boardPinNr uint8, railDeviceName string) (b *ButtonDevice, err error) {
	if err = boardsAPI.MapBinaryPin(boardID, boardPinNr, railDeviceName); err != nil {
		return
	}
	b = &ButtonDevice{
		name:      railDeviceName,
		boardsAPI: boardsAPI,
	}
	return
}

// StateChanged states true when Button status was changed
func (b *ButtonDevice) StateChanged() (hasChanged bool, err error) {
	var value uint8
	if value, err = b.boardsAPI.GetValue(b.name); err != nil {
		err = fmt.Errorf("Can't read value from '%s', %w", b.name, err)
		return
	}
	newState := value > 0
	if newState != b.oldState {
		b.oldState = newState
		hasChanged = true
	}
	return
}

// IsPressed gets the state of the button
func (b *ButtonDevice) IsPressed() bool {
	return b.oldState
}

// Name gets the name of the Button (rail device name)
func (b *ButtonDevice) Name() string {
	return b.name
}
