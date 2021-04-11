package raildevices

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommonOutputNew(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	// act
	co := NewCommonOutput("lamp dev", Timing{}, "lamp")
	// assert
	require.NotNil(co)
	assert.Equal("lamp dev", co.railDeviceName)
	assert.Equal("lamp", co.label)
	assert.Equal(false, co.IsOn())
	assert.Nil(co.IsDefective())
	stateChanged, _ := co.StateChanged("v")
	assert.Equal(true, stateChanged)
}

func TestCommonOutputIsOnSetStateStateChanged(t *testing.T) {
	// arrange
	assert := assert.New(t)
	co := NewCommonOutput("lamp dev", Timing{}, "lamp")
	// act
	stateChanged0, _ := co.StateChanged("v")
	state0 := co.IsOn()
	//
	co.SetState(true)
	stateChanged1, _ := co.StateChanged("v")
	state1 := co.IsOn()
	//
	co.SetState(true)
	stateChanged2, _ := co.StateChanged("v")
	state2 := co.IsOn()
	//
	co.SetState(false)
	stateChanged3, _ := co.StateChanged("v")
	state3 := co.IsOn()
	//
	co.SetState(false)
	stateChanged4, _ := co.StateChanged("v")
	state4 := co.IsOn()
	// assert
	assert.Equal(true, stateChanged0)
	assert.Equal(false, state0)
	assert.Equal(true, stateChanged1)
	assert.Equal(true, state1)
	assert.Equal(false, stateChanged2)
	assert.Equal(true, state2)
	assert.Equal(true, stateChanged3)
	assert.Equal(false, state3)
	assert.Equal(false, stateChanged4)
	assert.Equal(false, state4)
}

func TestCommonOutputIsDefectiveMakeDefectiveRepairStateChanged(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	co := NewCommonOutput("lamp dev", Timing{}, "lamp")
	// act
	deferr0 := co.IsDefective()
	stateChanged0, _ := co.StateChanged("v")
	//
	err1 := co.MakeDefective(func() error { return nil })
	stateChanged1, _ := co.StateChanged("v")
	deferr1 := co.IsDefective()
	//
	err2 := co.Repair()
	stateChanged2, _ := co.StateChanged("v")
	deferr2 := co.IsDefective()
	// assert
	require.Nil(deferr0)
	assert.Equal(true, stateChanged0)
	require.Nil(err1)
	assert.Equal(false, stateChanged1)
	assert.NotNil(deferr1)
	require.Nil(err2)
	assert.Equal(false, stateChanged2)
	assert.Nil(deferr2)
}

func TestCommonOutputMakeDefectiveWillSwitchOff(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	co := NewCommonOutput("lamp dev", Timing{}, "lamp")
	// act
	co.SetState(true)
	err := co.MakeDefective(func() error {
		co.SetState(false)
		return nil
	})
	stateChanged, _ := co.StateChanged("v")
	isOn := co.IsOn()
	// assert
	require.Nil(err)
	assert.Equal(false, isOn)
	assert.Equal(true, stateChanged)
}
