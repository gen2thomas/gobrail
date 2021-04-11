package raildevices

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToggleButtonNew(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	rm := ReadMock{}
	input := NewInputMock(&rm)
	// act
	toggleButton := NewToggleButton(input, "ToggleButton")
	// assert
	require.NotNil(toggleButton)
	assert.Equal("ToggleButton", toggleButton.railDeviceName)
	assert.Equal(0, rm.callCounter)
}

func TestToggleButtonToggleStateChangedIsOn(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	rm := ReadMock{values: [...]uint8{1, 1, 0, 0, 1}}
	input := NewInputMock(&rm)
	toggleButton := NewToggleButton(input, "ToggleButton")
	// act
	state0 := toggleButton.IsOn()
	changed1, err1 := toggleButton.StateChanged("v")
	state1 := toggleButton.IsOn()
	changed2, err2 := toggleButton.StateChanged("v")
	state2 := toggleButton.IsOn()
	changed3, err3 := toggleButton.StateChanged("v")
	state3 := toggleButton.IsOn()
	changed4, err4 := toggleButton.StateChanged("v")
	state4 := toggleButton.IsOn()
	changed5, err5 := toggleButton.StateChanged("v")
	state5 := toggleButton.IsOn()
	// assert
	require.Nil(err1)
	require.Nil(err2)
	require.Nil(err3)
	require.Nil(err4)
	require.Nil(err5)
	assert.Equal(5, rm.callCounter)
	assert.Equal(false, state0)
	assert.Equal(true, changed1)
	assert.Equal(true, state1)
	assert.Equal(false, changed2)
	assert.Equal(true, state2)
	assert.Equal(false, changed3)
	assert.Equal(true, state3)
	assert.Equal(false, changed4)
	assert.Equal(true, state4)
	assert.Equal(true, changed5)
	assert.Equal(false, state5)
}

func TestToggleButtonToggleStateChangedWhenReadErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	expectedError := fmt.Errorf("an error")
	rm := ReadMock{simError: expectedError}
	input := NewInputMock(&rm)
	toggleButton := NewToggleButton(input, "ToggleButton")
	// act
	_, err := toggleButton.StateChanged("v")
	// assert
	require.NotNil(err)
	assert.Equal(1, rm.callCounter)
	assert.Contains(err.Error(), "Can't read value from")
	assert.Equal(expectedError, errors.Unwrap(err))
}
