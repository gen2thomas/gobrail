package raildevices

// A Button is a rail device used for an input by button

import (
	"fmt"
)

// ButtonDevice is describes a Button
type ButtonDevice struct {
	name      string
	boardsAPI BoardsAPIer
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

// Name gets the name of the Button (rail device name)
func (l *ButtonDevice) Name() string {
	return l.name
}
