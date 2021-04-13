package raildevices

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTwoLightsSignalNew(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	co := NewCommonOutput("tls dev", Timing{})
	output1 := NewOutputMock(&WriteMock{})
	output2 := NewOutputMock(&WriteMock{})
	// act
	tls := NewTwoLightsSignal(co, output1, output2)
	// assert
	require.NotNil(tls)
	require.NotNil(tls.CommonOutputDevice)
	assert.Equal(co, tls.CommonOutputDevice)
	assert.Equal(output1, tls.outputPass)
	assert.Equal(output2, tls.outputStop)
}

func TestTwoLightsSignalSwitchOn(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	co := NewCommonOutput("tls dev", Timing{})
	wmPass := WriteMock{}
	wmStop := WriteMock{}
	outputPass := NewOutputMock(&wmPass)
	outputStop := NewOutputMock(&wmStop)
	tls := NewTwoLightsSignal(co, outputPass, outputStop)
	// act
	err := tls.SwitchOn()
	// assert
	require.Nil(err)
	assert.Equal(1, wmStop.callCounter)
	require.Equal(1, wmPass.callCounter)
	assert.Equal(uint8(0), wmStop.values[0])
	assert.Equal(uint8(1), wmPass.values[0])
	assert.Equal(true, tls.IsOn())
}

func TestTwoLightsSignalSwitchOff(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	co := NewCommonOutput("tls dev", Timing{})
	wmPass := WriteMock{}
	wmStop := WriteMock{}
	outputPass := NewOutputMock(&wmPass)
	outputStop := NewOutputMock(&wmStop)
	tls := NewTwoLightsSignal(co, outputPass, outputStop)
	// act
	err := tls.SwitchOff()
	// assert
	require.Nil(err)
	require.Equal(1, wmStop.callCounter)
	assert.Equal(1, wmPass.callCounter)
	assert.Equal(uint8(1), wmStop.values[0])
	assert.Equal(uint8(0), wmPass.values[0])
	assert.Equal(false, tls.IsOn())
}
