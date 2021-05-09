package raildevicesapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newRunableDevice(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	rm := &runnerMock{}
	// act
	rd := newRunableDevice(rm)
	// assert
	require.NotNil(rd)
	assert.Equal(rm, rd.Runner)
	assert.Equal(true, rd.firstRun)
}

func TestConnect(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	rd := runableDevice{Runner: runnerMock{name: "rdk"}}
	id := inputerMock{}
	// act
	err := rd.Connect(id, false)
	// assert
	require.Nil(err)
	assert.Equal(id, rd.connectedInput)
	assert.Equal(false, rd.inputInversion)
}

func TestConnectInverse(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	rd := runableDevice{Runner: runnerMock{name: "rdk"}}
	id := inputerMock{}
	// act
	err := rd.Connect(id, true)
	// assert
	require.Nil(err)
	assert.Equal(id, rd.connectedInput)
	assert.Equal(true, rd.inputInversion)
}

func TestConnectWhenAlreadyConnectedGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	id1 := inputerMock{}
	id2 := inputerMock{}
	rd := runableDevice{Runner: runnerMock{name: "rdk"}, connectedInput: id1}
	// act
	err := rd.Connect(id2, false)
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "is already connected")
	assert.Equal(id1, rd.connectedInput)
}

func TestConnectWhenSelfConnectGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	rd := runableDevice{Runner: runnerMock{name: "rdk"}}
	// act
	err := rd.Connect(rd, false)
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "Circular")
	assert.Nil(rd.connectedInput)
}

func TestRun(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	id := inputerMock{}
	rd := runableDevice{Runner: runnerMock{name: "rdk"}, connectedInput: id, firstRun: true}
	// act
	err := rd.Run()
	// assert
	require.Nil(err)
	assert.Equal(false, rd.firstRun)
}

func TestRunWithoutInputGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	rd := runableDevice{Runner: runnerMock{name: "rdk"}}
	// act
	err := rd.Run()
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "map to an input first")
}

func TestRunWhenStateChangedErrGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	id := inputerMock{simStateChangedErr: true}
	rd := runableDevice{Runner: runnerMock{name: "rdk"}, connectedInput: id}
	// act
	err := rd.Run()
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "state changed error")
}

func TestRunNoStateChanged(t *testing.T) {
	// arrange
	assert := assert.New(t)
	id := inputerMock{simStateChangedErr: false}
	rd := runableDevice{Runner: runnerMock{name: "rdk", simOffErr: true, simOnErr: true}, connectedInput: id}
	// act
	err := rd.Run()
	// assert
	assert.Nil(err)
}

func TestRunWhenSwitchOffErrGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	id := inputerMock{}
	rd := runableDevice{Runner: runnerMock{name: "rdk", simOffErr: true}, connectedInput: id, firstRun: true}
	// act
	err := rd.Run()
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "off error")
}

func TestRunWhenSwitchOnErrGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	id := inputerMock{isOn: true}
	rd := runableDevice{Runner: runnerMock{name: "rdk", simOnErr: true}, connectedInput: id, firstRun: true}
	// act
	err := rd.Run()
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "on error")
}

func TestRunWithInputInversionWhenSwitchOffErrGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	id := inputerMock{isOn: true}
	rd := runableDevice{Runner: runnerMock{name: "rdk", simOffErr: true}, connectedInput: id, inputInversion: true, firstRun: true}
	// act
	err := rd.Run()
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "off error")
}

func TestRunWithInputInversionWhenSwitchOnErrGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	id := inputerMock{}
	rd := runableDevice{Runner: runnerMock{name: "rdk", simOnErr: true}, connectedInput: id, inputInversion: true, firstRun: true}
	// act
	err := rd.Run()
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "on error")
}

func TestRunWithStateChangedWhenSwitchOffErrGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	id := inputerMock{stateChanged: true}
	rd := runableDevice{Runner: runnerMock{name: "rdk", simOffErr: true}, connectedInput: id}
	// act
	err := rd.Run()
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "off error")
}

func TestRunWithStateChangedWhenSwitchOnErrGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	id := inputerMock{stateChanged: true, isOn: true}
	rd := runableDevice{Runner: runnerMock{name: "rdk", simOnErr: true}, connectedInput: id}
	// act
	err := rd.Run()
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "on error")
}

func TestReleaseInput(t *testing.T) {
	// arrange
	assert := assert.New(t)
	rd := runableDevice{Runner: runnerMock{name: "rdk"}, connectedInput: inputerMock{}}
	// act
	rd.ReleaseInput()
	// assert
	assert.Nil(rd.connectedInput)
}
