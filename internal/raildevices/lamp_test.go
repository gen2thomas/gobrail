package raildevices

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLampNew(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	co := NewCommonOutput("Lamp", Timing{})
	wm := WriteMock{}
	output := NewOutputMock(&wm)
	// act
	lamp := NewLamp(co, output)
	// assert
	require.NotNil(lamp)
	require.NotNil(lamp.CommonOutputDevice)
	assert.Equal(co, lamp.CommonOutputDevice)
	assert.Equal(output, lamp.output)
	assert.Equal(0, wm.callCounter)
}

func TestLampSwitchOn(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	co := NewCommonOutput("Lamp", Timing{})
	wm := WriteMock{}
	output := NewOutputMock(&wm)
	lamp := NewLamp(co, output)
	// act
	err := lamp.SwitchOn()
	// assert
	require.Nil(err)
	require.Equal(1, wm.callCounter)
	assert.Equal(uint8(1), wm.values[0])
	assert.Equal(true, lamp.IsOn())
}

func TestLampSwitchOnWriteValueErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	co := NewCommonOutput("Lamp", Timing{})
	expErr := errors.New("an error")
	wm := WriteMock{simError: expErr}
	output := NewOutputMock(&wm)
	lamp := NewLamp(co, output)
	// act
	err := lamp.SwitchOn()
	// assert
	require.Equal(expErr, err)
	require.Equal(1, wm.callCounter)
	assert.Equal(uint8(1), wm.values[0])
	assert.Equal(false, lamp.IsOn())
}

func TestLampSwitchOnFailsWhenDefective(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	co := NewCommonOutput("Lamp", Timing{})
	wm := WriteMock{}
	output := NewOutputMock(&wm)
	lamp := NewLamp(co, output)
	lamp.MakeDefective()
	// act
	err := lamp.SwitchOn()
	// assert
	require.NotNil(err)
	assert.Equal(false, lamp.IsOn())
}

func TestLampSwitchOff(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	co := NewCommonOutput("Lamp", Timing{})
	wm := WriteMock{}
	output := NewOutputMock(&wm)
	lamp := NewLamp(co, output)
	// act
	err := lamp.SwitchOff()
	// assert
	require.Nil(err)
	assert.Equal(false, lamp.IsOn())
}

func TestLampSwitchOffWriteValueErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	co := NewCommonOutput("Lamp", Timing{})
	expErr := errors.New("an error")
	wm := WriteMock{simError: expErr}
	output := NewOutputMock(&wm)
	lamp := NewLamp(co, output)
	// act
	err := lamp.SwitchOff()
	// assert
	require.Equal(expErr, err)
	assert.Equal(false, lamp.IsOn())
}

func TestLampMakeDefective(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	co := NewCommonOutput("Lamp", Timing{})
	wm := WriteMock{}
	output := NewOutputMock(&wm)
	lamp := NewLamp(co, output)
	// act
	err := lamp.MakeDefective()
	// assert
	require.Nil(err)
	require.Equal(1, wm.callCounter)
	assert.Equal(uint8(0), wm.values[0])
	assert.Equal(false, lamp.IsOn())
}
