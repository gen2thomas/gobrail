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
	expectedIOPin := uint8(3)
	api := new(BoardsAPIMock)
	var usedBoardPinNrIOMap uint8
	callCounterIOMap := 0
	api.apiMapBinaryImpl = func(boardID string, boardPinNr uint8, railDeviceName string) (err error) {
		callCounterIOMap++
		usedBoardPinNrIOMap = boardPinNr
		return nil
	}
	callCounterAnaMap := 0
	api.apiMapAnalogImpl = func(boardID string, boardPinNr uint8, railDeviceName string) (err error) {
		callCounterAnaMap++
		return nil
	}
	callCounterMemMap := 0
	api.apiMapMemoryImpl = func(boardID string, boardPinNrOrNegative int, railDeviceName string) (err error) {
		callCounterMemMap++
		return nil
	}
	callCounterSetValue := 0
	api.apiSetValueImpl = func(railDeviceName string, value uint8) (err error) {
		callCounterSetValue++
		return nil
	}
	callCounterGetValue := 0
	api.apiGetValueImpl = func(railDeviceName string) (value uint8, err error) {
		callCounterGetValue++
		return 0, nil
	}
	// act
	toggleButton, err := NewToggleButton(api, "boardID", expectedIOPin, "ToggleButton")
	// assert
	require.Nil(err)
	require.NotNil(toggleButton)
	assert.Equal("ToggleButton", toggleButton.railDeviceName)
	assert.Equal(0, callCounterAnaMap)
	assert.Equal(1, callCounterIOMap)
	assert.Equal(expectedIOPin, usedBoardPinNrIOMap)
	assert.Equal(0, callCounterMemMap)
	assert.Equal(0, callCounterGetValue)
	assert.Equal(0, callCounterSetValue)
}

func TestToggleButtonNewWhenBinMapErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	expectedError := fmt.Errorf("an error")
	api := NewBoardsAPIMock()
	api.apiMapBinaryImpl = func(boardID string, boardPinNr uint8, railDeviceName string) (err error) {
		return expectedError
	}
	// act
	_, err := NewToggleButton(api, "boardID", 2, "ToggleButton")
	// assert
	require.NotNil(err)
	assert.Equal(expectedError, err)
}

func TestToggleButtonToggleStateChangedIsOn(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	api := NewBoardsAPIMock()
	var getValues = [...]uint8{1, 1, 0, 0, 1}
	callCounter := -1
	api.apiGetValueImpl = func(railDeviceName string) (value uint8, err error) {
		callCounter++
		return getValues[callCounter], nil
	}
	toggleButton, _ := NewToggleButton(api, "boardID", 3, "ToggleButton")
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
	api := NewBoardsAPIMock()
	api.apiGetValueImpl = func(railDeviceName string) (value uint8, err error) {
		return 0, expectedError
	}
	toggleButton, _ := NewToggleButton(api, "boardID", 1, "ToggleButton")
	// act
	_, err := toggleButton.StateChanged("v")
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "Can't read value from")
	assert.Equal(expectedError, errors.Unwrap(err))
}
