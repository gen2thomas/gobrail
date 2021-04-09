package raildevices

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTurnoutNew(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	expectedIOPins := [...]uint8{4, 3}
	api := new(BoardsAPIMock)
	var usedBoardPinNrBinMaps [2]uint8
	api.apiMapBinaryImpl = func(boardID string, boardPinNr uint8, railDeviceName string) (err error) {
		usedBoardPinNrBinMaps[api.callCounterBinMap-1] = boardPinNr
		return nil
	}
	// act
	turnout, err := NewTurnout(api, "boardID", expectedIOPins[0], "turnout dev", expectedIOPins[1], Timing{})
	// assert
	require.Nil(err)
	require.NotNil(turnout)
	require.NotNil(turnout.cmnOutDev)
	assert.Equal("turnout dev branch", turnout.railDeviceNameBranch)
	assert.Equal(0, api.callCounterAnaMap)
	assert.Equal(2, api.callCounterBinMap)
	assert.Equal(0, api.callCounterMemMap)
	assert.Equal(0, api.callCounterGetValue)
	assert.Equal(0, api.callCounterSetValue)
	assert.Equal(expectedIOPins, usedBoardPinNrBinMaps)
}

func TestTurnoutNewWhenBinMapErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	expectedError := fmt.Errorf("an error")
	api := NewBoardsAPIMock()
	api.apiMapBinaryImpl = func(boardID string, boardPinNr uint8, railDeviceName string) (err error) {
		if api.callCounterBinMap == 2 {
			return expectedError
		}
		return
	}
	// act
	_, err := NewTurnout(api, "boardID", 2, "turnout dev", 3, Timing{})
	// assert
	require.NotNil(err)
	assert.Equal(expectedError, err)
}

func TestTurnoutSwitchOn(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	api := NewBoardsAPIMock()
	api.apiSetValueImpl = func(railDeviceName string, value uint8) (err error) {
		key := fmt.Sprintf("%s%d", railDeviceName, api.callCounterSetValue)
		api.values[key] = value
		return
	}
	turnout, _ := NewTurnout(api, "boardID", 1, "turnout dev", 2, Timing{})
	// act
	err := turnout.SwitchOn()
	// assert
	require.Nil(err)
	assert.Equal(2, api.callCounterSetValue)
	assert.Equal(uint8(1), api.values["turnout dev branch1"])
	assert.Equal(uint8(0), api.values["turnout dev branch2"])
	assert.Equal(true, turnout.IsOn())
}

func TestTurnoutSwitchOff(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	api := NewBoardsAPIMock()
	api.apiSetValueImpl = func(railDeviceName string, value uint8) (err error) {
		key := fmt.Sprintf("%s%d", railDeviceName, value)
		api.values[key] = value
		return
	}
	turnout, _ := NewTurnout(api, "boardID", 1, "turnout dev", 2, Timing{})
	// act
	err := turnout.SwitchOff()
	// assert
	require.Nil(err)
	assert.Equal(2, api.callCounterSetValue)
	assert.Equal(uint8(1), api.values["turnout dev1"])
	assert.Equal(uint8(0), api.values["turnout dev2"])
	assert.Equal(false, turnout.IsOn())
}
