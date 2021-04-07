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
	var usedBoardPinNrMemoryMap [2]int
	var usedSetValue [3]uint8
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
		usedBoardPinNrMemoryMap[callCounterMemMap] = boardPinNrOrNegative
		callCounterMemMap++
		return nil
	}
	callCounterSetValue := 0
	api.apiSetValueImpl = func(railDeviceName string, value uint8) (err error) {
		usedSetValue[callCounterSetValue] = value
		callCounterSetValue++
		return nil
	}
	callCounterGetValue := 0
	api.apiGetValueImpl = func(railDeviceName string) (value uint8, err error) {
		callCounterGetValue++
		return 0, nil
	}
	// act
	lamp, err := NewLamp(api, "boardID", expectedIOPin, "lamp", Timing{})
	// assert
	require.Nil(err)
	require.NotNil(lamp)
	assert.Equal("lamp", lamp.name)
	assert.Contains(lamp.stateName, "lamp")
	assert.Contains(lamp.defectiveName, "lamp")
	assert.Equal(0, callCounterAnaMap)
	assert.Equal(1, callCounterIOMap)
	assert.Equal(expectedIOPin, usedBoardPinNrIOMap)
	assert.Equal(2, callCounterMemMap)
	assert.Equal(-1, usedBoardPinNrMemoryMap[0])
	assert.Equal(-1, usedBoardPinNrMemoryMap[1])
	assert.Equal(1, callCounterGetValue)
	assert.Equal(3, callCounterSetValue)
	assert.Equal(uint8(0), usedSetValue[0])
	assert.Equal(uint8(0), usedSetValue[1])
	assert.Equal(uint8(0), usedSetValue[2])
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

func TestNewLampWhenFirstMemMapErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	expectedError := fmt.Errorf("an error")
	api := NewBoardsAPIMock()
	callCounter := 0
	api.apiMapMemoryImpl = func(boardID string, boardPinNrOrNegative int, railDeviceName string) (err error) {
		if callCounter == 0 {
			err = expectedError
		}
		callCounter++
		return
	}
	// act
	_, err := NewLamp(api, "boardID", 2, "lamp", Timing{})
	// assert
	require.NotNil(err)
	assert.Equal(expectedError, err)
}

func TestNewLampWhenSecondMemMapErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	expectedError := fmt.Errorf("an error")
	api := NewBoardsAPIMock()
	callCounter := 0
	api.apiMapMemoryImpl = func(boardID string, boardPinNrOrNegative int, railDeviceName string) (err error) {
		if callCounter == 1 {
			err = expectedError
		}
		callCounter++
		return
	}
	// act
	_, err := NewLamp(api, "boardID", 2, "lamp", Timing{})
	// assert
	require.NotNil(err)
	assert.Equal(expectedError, err)
}

func TestIsOnSwitchOnSwitchOff(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	api := NewBoardsAPIMock()
	lamp, _ := NewLamp(api, "boardID", 1, "lamp", Timing{})
	// act
	shouldBeNotOn1, err1 := lamp.IsOn()
	err2 := lamp.SwitchOn()
	shouldBeOn, err3 := lamp.IsOn()
	err4 := lamp.SwitchOff()
	shouldBeNotOn2, err5 := lamp.IsOn()
	// assert
	require.Nil(err1)
	require.Nil(err2)
	require.Nil(err3)
	require.Nil(err4)
	require.Nil(err5)
	assert.Equal(false, shouldBeNotOn1)
	assert.Equal(true, shouldBeOn)
	assert.Equal(false, shouldBeNotOn2)
}

func TestIsOnWhenGetValueHasErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	expectedError := fmt.Errorf("an error")
	api := NewBoardsAPIMock()
	callCounter := -1
	api.apiGetValueImpl = func(railDeviceName string) (value uint8, err error) {
		// first call comes from "NewLamp"
		callCounter++
		if callCounter == 1 {
			err = expectedError
		}

		return 1, err
	}
	lamp, _ := NewLamp(api, "boardID", 1, "lamp", Timing{})
	// act
	shouldBeNotOn, err := lamp.IsOn()
	// assert
	require.NotNil(err)
	assert.Equal(false, shouldBeNotOn)
}

func TestIsDefectiveMakeDefectiveRepair(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	api := NewBoardsAPIMock()
	lamp, _ := NewLamp(api, "boardID", 1, "lamp", Timing{})
	// act
	shouldBeNot1, err1 := lamp.IsDefective()
	err2 := lamp.MakeDefective()
	shouldBe, err3 := lamp.IsDefective()
	err4 := lamp.Repair()
	shouldBeNot2, err5 := lamp.IsDefective()
	// assert
	require.Nil(err1)
	require.Nil(err2)
	require.Nil(err3)
	require.Nil(err4)
	require.Nil(err5)
	assert.Equal(false, shouldBeNot1)
	assert.Equal(true, shouldBe)
	assert.Equal(false, shouldBeNot2)
}

func TestIsDefectiveWhenGetValueHasErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	expectedError := fmt.Errorf("an error")
	api := NewBoardsAPIMock()
	callCounter := -1
	api.apiGetValueImpl = func(railDeviceName string) (value uint8, err error) {
		// first call comes from "NewLamp"
		callCounter++
		if callCounter == 1 {
			err = expectedError
		}

		return 1, err
	}
	lamp, _ := NewLamp(api, "boardID", 1, "lamp", Timing{})
	// act
	shouldBeNot, err := lamp.IsDefective()
	// assert
	require.NotNil(err)
	assert.Equal(false, shouldBeNot)
}

func TestSwitchOnWhenIsDefectiveGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	api := NewBoardsAPIMock()
	lamp, _ := NewLamp(api, "boardID", 1, "lamp", Timing{})
	// act
	err1 := lamp.MakeDefective()
	err2 := lamp.SwitchOn()
	// assert
	require.Nil(err1)
	assert.NotNil(err2)
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
	isOn, err3 := lamp.IsOn()
	// assert
	require.Nil(err1)
	require.Nil(err2)
	require.Nil(err3)
	assert.Equal(false, isOn)
}
