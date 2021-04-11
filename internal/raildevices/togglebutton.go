package raildevices

// A ToggleButton is a rail device used for an input by a button
// the output will change on each press of button

import (
	"fmt"

	"github.com/gen2thomas/gobrail/internal/boardpin"
)

// ToggleButtonDevice describes a ToggleButton
type ToggleButtonDevice struct {
	railDeviceName string
	oldState       bool
	toggleState    bool
	oldToggleState map[string]bool
	input          *boardpin.Input
}

// NewToggleButton creates an instance of a ToggleButton
func NewToggleButton(input *boardpin.Input, railDeviceName string) (ld *ToggleButtonDevice) {
	ld = &ToggleButtonDevice{
		railDeviceName: railDeviceName,
		oldToggleState: make(map[string]bool),
		input:          input,
	}
	return
}

// StateChanged states true when ToggleButton status was changed
func (b *ToggleButtonDevice) StateChanged(visitor string) (hasChanged bool, err error) {
	var value uint8
	if value, err = b.input.ReadValue(); err != nil {
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
