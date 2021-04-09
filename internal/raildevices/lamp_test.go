package raildevices

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLamp(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	expectedIOPin := uint8(5)
	api := new(BoardsAPIMock)
	var usedBoardPinNrIOMap uint8
	api.apiMapBinaryImpl = func(boardID string, boardPinNr uint8, railDeviceName string) (err error) {
		usedBoardPinNrIOMap = boardPinNr
		return nil
	}
	// act
	lamp, err := NewLamp(api, "boardID", expectedIOPin, "lamp", Timing{})
	// assert
	require.Nil(err)
	require.NotNil(lamp)
	assert.Equal("lamp", lamp.name)
	assert.Equal(false, lamp.IsOn())
	assert.Equal(false, lamp.IsDefective())
	stateChanged, _ := lamp.StateChanged("v")
	assert.Equal(true, stateChanged)
	assert.Equal(0, api.callCounterAnaMap)
	assert.Equal(1, api.callCounterBinMap)
	assert.Equal(0, api.callCounterMemMap)
	assert.Equal(0, api.callCounterGetValue)
	assert.Equal(0, api.callCounterSetValue)
	assert.Equal(expectedIOPin, usedBoardPinNrIOMap)
}

func TestNewLampWhenBinMapErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	expectedError := fmt.Errorf("an error")
	api := NewBoardsAPIMock()
	api.apiMapBinaryImpl = func(boardID string, boardPinNr uint8, railDeviceName string) (err error) {
		return expectedError
	}
	// act
	_, err := NewLamp(api, "boardID", 2, "lamp", Timing{})
	// assert
	require.NotNil(err)
	assert.Equal(expectedError, err)
}

func TestIsOnSwitchOnSwitchOffStateChanged(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	api := NewBoardsAPIMock()
	lamp, _ := NewLamp(api, "boardID", 1, "lamp", Timing{})
	// act
	stateChanged0, _ := lamp.StateChanged("v")
	state0 := lamp.IsOn()
	//
	err1 := lamp.SwitchOn()
	stateChanged1, _ := lamp.StateChanged("v")
	state1 := lamp.IsOn()
	//
	err2 := lamp.SwitchOn()
	stateChanged2, _ := lamp.StateChanged("v")
	state2 := lamp.IsOn()
	//
	err3 := lamp.SwitchOff()
	stateChanged3, _ := lamp.StateChanged("v")
	state3 := lamp.IsOn()
	//
	err4 := lamp.SwitchOff()
	stateChanged4, _ := lamp.StateChanged("v")
	state4 := lamp.IsOn()
	// assert
	assert.Equal(true, stateChanged0)
	assert.Equal(false, state0)
	require.Nil(err1)
	assert.Equal(true, stateChanged1)
	assert.Equal(true, state1)
	require.Nil(err2)
	assert.Equal(false, stateChanged2)
	assert.Equal(true, state2)
	require.Nil(err3)
	assert.Equal(true, stateChanged3)
	assert.Equal(false, state3)
	require.Nil(err4)
	assert.Equal(false, stateChanged4)
	assert.Equal(false, state4)
}

func TestIsDefectiveMakeDefectiveRepairStateChanged(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	api := NewBoardsAPIMock()
	lamp, _ := NewLamp(api, "boardID", 1, "lamp", Timing{})
	// act
	state0 := lamp.IsDefective()
	stateChanged0, _ := lamp.StateChanged("v")
	//
	err1 := lamp.MakeDefective()
	stateChanged1, _ := lamp.StateChanged("v")
	state1 := lamp.IsDefective()
	//
	err2 := lamp.Repair()
	stateChanged2, _ := lamp.StateChanged("v")
	state2 := lamp.IsDefective()
	// assert
	require.Nil(err1)
	require.Nil(err2)
	assert.Equal(true, stateChanged0)
	assert.Equal(false, state0)
	assert.Equal(false, stateChanged1)
	assert.Equal(true, state1)
	assert.Equal(false, stateChanged2)
	assert.Equal(false, state2)
}

func TestSwitchOnWhenIsDefectiveGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	api := NewBoardsAPIMock()
	lamp, _ := NewLamp(api, "boardID", 1, "lamp", Timing{})
	// act
	stateChanged0, _ := lamp.StateChanged("v")
	err1 := lamp.MakeDefective()
	err2 := lamp.SwitchOn()
	stateChanged1, _ := lamp.StateChanged("v")
	// assert
	require.Nil(err1)
	assert.NotNil(err2)
	assert.Equal(true, stateChanged0)
	assert.Equal(false, stateChanged1)
}

func TestMakeDefectiveWillSwitchOff(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	api := NewBoardsAPIMock()
	lamp, _ := NewLamp(api, "boardID", 1, "lamp", Timing{})
	// act
	err1 := lamp.SwitchOn()
	err2 := lamp.MakeDefective()
	stateChanged, _ := lamp.StateChanged("v")
	isOn := lamp.IsOn()
	// assert
	require.Nil(err1)
	require.Nil(err2)
	assert.Equal(false, isOn)
	assert.Equal(true, stateChanged)
}
