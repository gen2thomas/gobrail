package raildevices

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTwoLightSignalNew(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	expectedIOPins := [...]uint8{5, 7}
	api := new(BoardsAPIMock)
	var usedBoardPinNrBinMaps [2]uint8
	api.apiMapBinaryImpl = func(boardID string, boardPinNr uint8, railDeviceName string) (err error) {
		usedBoardPinNrBinMaps[api.callCounterBinMap-1] = boardPinNr
		return nil
	}
	// act
	tls, err := NewTwoLightSignal(api, "boardID", expectedIOPins[0], "tls dev", expectedIOPins[1], Timing{})
	// assert
	require.Nil(err)
	require.NotNil(tls)
	require.NotNil(tls.cmnOutDev)
	assert.Equal("tls dev stop", tls.railDeviceNameStop)
	assert.Equal(0, api.callCounterAnaMap)
	assert.Equal(2, api.callCounterBinMap)
	assert.Equal(0, api.callCounterMemMap)
	assert.Equal(0, api.callCounterGetValue)
	assert.Equal(0, api.callCounterSetValue)
	assert.Equal(expectedIOPins, usedBoardPinNrBinMaps)
}

func TestTwoLightSignalNewWhenBinMapErrorGetsError(t *testing.T) {
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
	_, err := NewTwoLightSignal(api, "boardID", 2, "tls dev", 3, Timing{})
	// assert
	require.NotNil(err)
	assert.Equal(expectedError, err)
}

func TestTwoLightSignalSwitchOn(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	api := NewBoardsAPIMock()
	tls, _ := NewTwoLightSignal(api, "boardID", 1, "tls dev", 2, Timing{})
	// act
	err := tls.SwitchOn()
	// assert
	require.Nil(err)
	assert.Equal(2, api.callCounterSetValue)
	assert.Equal(uint8(1), api.values["tls dev"])
	assert.Equal(uint8(0), api.values["tls dev stop"])
	assert.Equal(true, tls.IsOn())
}

func TestTwoLightSignalSwitchOff(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	api := NewBoardsAPIMock()
	tls, _ := NewTwoLightSignal(api, "boardID", 1, "tls dev", 2, Timing{})
	// act
	err := tls.SwitchOff()
	// assert
	require.Nil(err)
	assert.Equal(2, api.callCounterSetValue)
	assert.Equal(uint8(0), api.values["tls dev"])
	assert.Equal(uint8(1), api.values["tls dev stop"])
	assert.Equal(false, tls.IsOn())
}
