package raildevices

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTurnoutNew(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	co := NewCommonOutput("turnout dev", Timing{})
	output1 := NewOutputMock(&WriteMock{})
	output2 := NewOutputMock(&WriteMock{})
	// act
	turnout := NewTurnout(co, output1, output2)
	// assert
	require.NotNil(turnout)
	require.NotNil(turnout.CommonOutputDevice)
	assert.Equal(co, turnout.CommonOutputDevice)
	assert.Equal(output1, turnout.outputBranch)
	assert.Equal(output2, turnout.outputMain)
}

func TestTurnoutSwitchOn(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	co := NewCommonOutput("turnout dev", Timing{})
	wmBranch := WriteMock{}
	wmMain := WriteMock{}
	outputBranch := NewOutputMock(&wmBranch)
	outputMain := NewOutputMock(&wmMain)
	turnout := NewTurnout(co, outputBranch, outputMain)
	// act
	err := turnout.SwitchOn()
	// assert
	require.Nil(err)
	require.Equal(2, wmBranch.callCounter)
	assert.Equal(0, wmMain.callCounter)
	assert.Equal(uint8(1), wmBranch.values[0])
	assert.Equal(uint8(0), wmBranch.values[1])
	assert.Equal(true, turnout.IsOn())
}

func TestTurnoutSwitchOnWhenErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	co := NewCommonOutput("turnout dev", Timing{})
	expErr := errors.New("an error")
	wmBranch := WriteMock{simError: expErr}
	wmMain := WriteMock{}
	outputBranch := NewOutputMock(&wmBranch)
	outputMain := NewOutputMock(&wmMain)
	turnout := NewTurnout(co, outputBranch, outputMain)
	// act
	err := turnout.SwitchOn()
	// assert
	require.NotNil(err)
	assert.Equal(expErr, err)
	assert.Equal(false, turnout.IsOn())
}

func TestTurnoutSwitchOff(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	co := NewCommonOutput("turnout dev", Timing{})
	wmBranch := WriteMock{}
	wmMain := WriteMock{}
	outputBranch := NewOutputMock(&wmBranch)
	outputMain := NewOutputMock(&wmMain)
	turnout := NewTurnout(co, outputBranch, outputMain)
	// act
	err := turnout.SwitchOff()
	// assert
	require.Nil(err)
	require.Equal(2, wmMain.callCounter)
	assert.Equal(0, wmBranch.callCounter)
	assert.Equal(uint8(1), wmMain.values[0])
	assert.Equal(uint8(0), wmMain.values[1])
	assert.Equal(false, turnout.IsOn())
}

func TestTurnoutSwitchOffWhenErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	co := NewCommonOutput("turnout dev", Timing{})
	expErr := errors.New("an error")
	wmBranch := WriteMock{}
	wmMain := WriteMock{simError: expErr}
	outputBranch := NewOutputMock(&wmBranch)
	outputMain := NewOutputMock(&wmMain)
	turnout := NewTurnout(co, outputBranch, outputMain)
	// act
	err := turnout.SwitchOff()
	// assert
	require.NotNil(err)
	assert.Equal(expErr, err)
	assert.Equal(false, turnout.IsOn())
}
