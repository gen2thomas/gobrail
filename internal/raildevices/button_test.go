package raildevices

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewButton(t *testing.T) {
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
	button, err := NewButton(api, "boardID", expectedIOPin, "Button")
	// assert
	require.Nil(err)
	require.NotNil(button)
	assert.Equal("Button", button.name)
	assert.Equal(0, callCounterAnaMap)
	assert.Equal(1, callCounterIOMap)
	assert.Equal(expectedIOPin, usedBoardPinNrIOMap)
	assert.Equal(0, callCounterMemMap)
	assert.Equal(0, callCounterGetValue)
	assert.Equal(0, callCounterSetValue)
}

func TestNewButtonWhenBinMapErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	expectedError := fmt.Errorf("an error")
	api := NewBoardsAPIMock()
	api.apiMapBinaryImpl = func(boardID string, boardPinNr uint8, railDeviceName string) (err error) {
		return expectedError
	}
	// act
	_, err := NewButton(api, "boardID", 2, "Button")
	// assert
	require.NotNil(err)
	assert.Equal(expectedError, err)
}

func TestStateChangedIsPressed(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	api := NewBoardsAPIMock()
	var getValues = [...]uint8{1, 1, 0, 0}
	callCounter := -1
	api.apiGetValueImpl = func(railDeviceName string) (value uint8, err error) {
		callCounter++
		return getValues[callCounter], nil
	}
	button, _ := NewButton(api, "boardID", 3, "Button")
	// act
	state0 := button.IsPressed()
	changed1, err1 := button.StateChanged()
	state1 := button.IsPressed()
	changed2, err2 := button.StateChanged()
	state2 := button.IsPressed()
	changed3, err3 := button.StateChanged()
	state3 := button.IsPressed()
	changed4, err4 := button.StateChanged()
	state4 := button.IsPressed()
	// assert
	require.Nil(err1)
	require.Nil(err2)
	require.Nil(err3)
	require.Nil(err4)
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
	api := NewBoardsAPIMock()
	api.apiGetValueImpl = func(railDeviceName string) (value uint8, err error) {
		return 0, expectedError
	}
	button, _ := NewButton(api, "boardID", 1, "Button")
	// act
	_, err := button.StateChanged()
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "Can't read value from")
	assert.Equal(expectedError, errors.Unwrap(err))
}
