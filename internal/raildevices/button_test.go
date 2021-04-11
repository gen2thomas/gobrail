package raildevices

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestButtonNew(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	input := NewInputMock(&ReadMock{})
	// act
	button := NewButton(input, "Button")
	// assert
	require.NotNil(button)
	assert.Equal("Button", button.railDeviceName)
	assert.Equal(input, button.input)
}

func TestButtonStateChangedIsOn(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	rm := ReadMock{values: [...]uint8{1, 1, 0, 0, 255}}
	input := NewInputMock(&rm)
	button := NewButton(input, "Button")
	// act
	state0 := button.IsOn()
	changed1, err1 := button.StateChanged("v")
	state1 := button.IsOn()
	changed2, err2 := button.StateChanged("v")
	state2 := button.IsOn()
	changed3, err3 := button.StateChanged("v")
	state3 := button.IsOn()
	changed4, err4 := button.StateChanged("v")
	state4 := button.IsOn()
	// assert
	require.Nil(err1)
	require.Nil(err2)
	require.Nil(err3)
	require.Nil(err4)
	require.Equal(4, rm.callCounter)
	assert.Equal(false, state0)

	assert.Equal(true, changed1)
	assert.Equal(true, state1)

	assert.Equal(false, changed2)
	assert.Equal(true, state2)

	assert.Equal(true, changed3)
	assert.Equal(false, state3)

	assert.Equal(false, changed4)
	assert.Equal(false, state4)
}

func TestButtonStateChangedWhenReadErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	expectedError := fmt.Errorf("an error")
	rm := ReadMock{simError: expectedError}
	input := NewInputMock(&rm)
	button := NewButton(input, "Button")
	// act
	_, err := button.StateChanged("v")
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "Can't read value from")
	assert.Equal(expectedError, errors.Unwrap(err))
}
