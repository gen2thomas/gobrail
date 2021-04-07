package raildevices

// A Button is a rail device used for an input by button

import (
	"fmt"
)

// ButtonDevice is describes a Button
type ButtonDevice struct {
	name       string
	wasPressed bool
	boardsAPI  BoardsAPIer
}

// NewButton creates an instance of a Button
func NewButton(boardsAPI BoardsAPIer, boardID string, boardPinNr uint8, railDeviceName string) (ld *ButtonDevice, err error) {
	if err = boardsAPI.MapBinaryPin(boardID, boardPinNr, railDeviceName); err != nil {
		return
	}
	ld = &ButtonDevice{
		name:      railDeviceName,
		boardsAPI: boardsAPI,
	}
	return
}

// IsPressed states true when Button is pressed
func (l *ButtonDevice) IsPressed() (isPressed bool, err error) {
	var value uint8
	if value, err = l.boardsAPI.GetValue(l.name); err != nil {
		err = fmt.Errorf("Can't read value from '%s', %w", l.name, err)
		return
	}
	return value > 0, nil
}

// IsChanged states true when Button status was changed
// The value can be read by WasPressed()
func (l *ButtonDevice) IsChanged() (isChanged bool, err error) {
	var value uint8
	if value, err = l.boardsAPI.GetValue(l.name); err != nil {
		err = fmt.Errorf("Can't read value from '%s', %w", l.name, err)
		return
	}
	isPressed := value > 0
	if isPressed == l.wasPressed {
		return
	}
	l.wasPressed = isPressed
	return true, nil
}

// WasPressed gets the state of the button last read from input
// means last call of IsPressed() or IsChanged()
func (l *ButtonDevice) WasPressed() bool {
	return l.wasPressed
}

// Name gets the name of the Button (rail device name)
func (l *ButtonDevice) Name() string {
	return l.name
}
