package raildevices

// A Button is a rail device used for an input by button

import (
	"fmt"

	"github.com/gen2thomas/gobrail/internal/boardpin"
)

// ButtonDevice describes a Button
type ButtonDevice struct {
	railDeviceName string
	state          bool
	oldState       map[string]bool
	input          *boardpin.Input
}

// NewButton creates an instance of a Button
func NewButton(input *boardpin.Input, railDeviceName string) (b *ButtonDevice) {
	b = &ButtonDevice{
		railDeviceName: railDeviceName,
		oldState:       make(map[string]bool),
		input:          input,
	}
	return
}

// StateChanged states true when Button status was changed
func (b *ButtonDevice) StateChanged(visitor string) (hasChanged bool, err error) {
	var value uint8
	if value, err = b.input.ReadValue(); err != nil {
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
