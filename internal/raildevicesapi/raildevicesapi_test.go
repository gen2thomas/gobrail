package raildevicesapi

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gen2thomas/gobrail/internal/boardpin"
	"github.com/gen2thomas/gobrail/internal/devicerecipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type boardsIOAPIMock struct{}

type inputerMock struct{}
type runnerMock struct{ name string }

func TestNewRailDevicesAPI(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	ba := &boardsIOAPIMock{}
	// act
	da := NewRailDevicesAPI(ba)
	// assert
	require.NotNil(da)
	assert.NotNil(da.devices)
	assert.NotNil(da.runableDevices)
	assert.NotNil(da.inputDevices)
	assert.NotNil(da.connections)
	assert.Equal(da.boardsIOAPI, ba)
}

func TestAddDeviceExistGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	da := RailDeviceAPI{}
	da.devices = map[string]struct{}{"test_device": struct{}{}}
	// act
	err := da.AddDevice(devicerecipe.Ingredients{Name: "test_device"})
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "already in use")
}

func TestAddDeviceUnknownGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	da := RailDeviceAPI{}
	da.devices = make(map[string]struct{})
	// act
	err := da.AddDevice(devicerecipe.Ingredients{Name: "test_device", Type: "TypUnknown"})
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "Unknown type")
}

func TestAddDevice(t *testing.T) {
	var addDeviceTests = map[string]devicerecipe.Ingredients{
		"AddButton":          {Name: "test_device", Type: "Button", BoardID: "test_board", BoardPinNrPrim: 0},
		"AddToggleButton":    {Name: "test_device", Type: "ToggleButton", BoardID: "test_board", BoardPinNrPrim: 1},
		"AddLamp":            {Name: "test_device", Type: "Lamp", BoardID: "test_board", BoardPinNrPrim: 2},
		"AddTwoLightsSignal": {Name: "test_device", Type: "TwoLightsSignal", BoardID: "test_board", BoardPinNrPrim: 3, BoardPinNrSec: 4},
		"AddTurnout":         {Name: "test_device", Type: "Turnout", BoardID: "test_board", BoardPinNrPrim: 5, BoardPinNrSec: 6, Connect: "test_connect"},
	}
	for name, at := range addDeviceTests {
		t.Run(name, func(t *testing.T) {
			// arrange
			assert := assert.New(t)
			require := require.New(t)
			ba := &boardsIOAPIMock{}
			da := RailDeviceAPI{boardsIOAPI: ba}
			da.devices = make(map[string]struct{})
			da.inputDevices = make(map[string]Inputer)
			da.runableDevices = make(map[string]*runableDevice)
			da.connections = make(map[string]connection)
			// act
			err := da.AddDevice(at)
			// assert
			require.Nil(err)
			assert.Contains(da.devices, "test_device")
			if strings.Contains(name, "Button") {
				assert.Contains(da.inputDevices, "test_device")
				assert.NotContains(da.runableDevices, "test_device")
			} else {
				assert.Contains(da.runableDevices, "test_device")
				assert.NotContains(da.inputDevices, "test_device")
			}
			if at.Connect != "" {
				assert.Equal("test_connect", da.connections["test_device"].name)
			} else {
				assert.NotContains(da.connections, "test_device")
			}
		})
	}
}

func TestConnectNow(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	da := RailDeviceAPI{}
	// inp_dev --> in_run_dev --> run_dev
	da.inputDevices = map[string]Inputer{"inp_dev_key": &inputerMock{}}
	da.runableDevices = map[string]*runableDevice{"run_dev_key": &runableDevice{Runner: runnerMock{name: "rdk"}}, "in_run_dev_key": &runableDevice{Runner: runnerMock{name: "irdk"}}}
	da.connections = map[string]connection{"run_dev_key": connection{name: "in_run_dev_key"}, "in_run_dev_key": connection{name: "inp_dev_key"}}
	// act
	err := da.ConnectNow()
	// assert
	require.Nil(err)
	assert.Equal(da.inputDevices["inp_dev_key"], da.runableDevices["in_run_dev_key"].connectedInput)
	assert.Equal(da.runableDevices["in_run_dev_key"], da.runableDevices["run_dev_key"].connectedInput)
}

func Test_createButton(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	ba := &boardsIOAPIMock{}
	da := RailDeviceAPI{boardsIOAPI: ba}
	da.devices = make(map[string]struct{})
	// act
	inp, err := da.createButton(devicerecipe.Ingredients{})
	// assert
	require.Nil(err)
	assert.NotNil(inp)
}

func Test_createButtonGetInputPinErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	ba := &boardsIOAPIMock{}
	da := RailDeviceAPI{boardsIOAPI: ba}
	da.devices = make(map[string]struct{})
	// act
	_, err := da.createButton(devicerecipe.Ingredients{BoardID: "error"})
	// assert
	require.NotNil(err)
	assert.Equal("test error", err.Error())
}

func Test_createToggleButton(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	ba := &boardsIOAPIMock{}
	da := RailDeviceAPI{boardsIOAPI: ba}
	da.devices = make(map[string]struct{})
	// act
	inp, err := da.createToggleButton(devicerecipe.Ingredients{})
	// assert
	require.Nil(err)
	assert.NotNil(inp)
}

func Test_createToggleButtonGetInputPinErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	ba := &boardsIOAPIMock{}
	da := RailDeviceAPI{boardsIOAPI: ba}
	da.devices = make(map[string]struct{})
	// act
	_, err := da.createToggleButton(devicerecipe.Ingredients{BoardID: "error"})
	// assert
	require.NotNil(err)
	assert.Equal("test error", err.Error())
}

func Test_createLamp(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	ba := &boardsIOAPIMock{}
	da := RailDeviceAPI{boardsIOAPI: ba}
	da.devices = make(map[string]struct{})
	// act
	outp, err := da.createLamp(devicerecipe.Ingredients{})
	// assert
	require.Nil(err)
	assert.NotNil(outp)
}

func Test_createLampGetOutputPinErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	ba := &boardsIOAPIMock{}
	da := RailDeviceAPI{boardsIOAPI: ba}
	da.devices = make(map[string]struct{})
	// act
	_, err := da.createLamp(devicerecipe.Ingredients{BoardID: "error", BoardPinNrPrim: 88})
	// assert
	require.NotNil(err)
	assert.Equal("test error", err.Error())
}

func Test_createTwoLightSignal(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	ba := &boardsIOAPIMock{}
	da := RailDeviceAPI{boardsIOAPI: ba}
	da.devices = make(map[string]struct{})
	// act
	outp, err := da.createTwoLightSignal(devicerecipe.Ingredients{})
	// assert
	require.Nil(err)
	assert.NotNil(outp)
}

func Test_createTwoLightSignalGetOutPinPrimErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	ba := &boardsIOAPIMock{}
	da := RailDeviceAPI{boardsIOAPI: ba}
	da.devices = make(map[string]struct{})
	// act
	_, err := da.createTwoLightSignal(devicerecipe.Ingredients{BoardID: "error", BoardPinNrPrim: 88})
	// assert
	require.NotNil(err)
	assert.Equal("test error", err.Error())
}

func Test_createTwoLightSignalGetOutPinSecErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	ba := &boardsIOAPIMock{}
	da := RailDeviceAPI{boardsIOAPI: ba}
	da.devices = make(map[string]struct{})
	// act
	_, err := da.createTwoLightSignal(devicerecipe.Ingredients{BoardID: "error", BoardPinNrSec: 88})
	// assert
	require.NotNil(err)
	assert.Equal("test error", err.Error())
}

func Test_createTurnout(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	ba := &boardsIOAPIMock{}
	da := RailDeviceAPI{boardsIOAPI: ba}
	da.devices = make(map[string]struct{})
	// act
	outp, err := da.createTurnout(devicerecipe.Ingredients{})
	// assert
	require.Nil(err)
	assert.NotNil(outp)
}

func Test_createTurnoutGetOutPinPrimErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	ba := &boardsIOAPIMock{}
	da := RailDeviceAPI{boardsIOAPI: ba}
	da.devices = make(map[string]struct{})
	// act
	_, err := da.createTurnout(devicerecipe.Ingredients{BoardID: "error", BoardPinNrPrim: 88})
	// assert
	require.NotNil(err)
	assert.Equal("test error", err.Error())
}

func Test_createTurnoutGetOutPinSecErrorGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	ba := &boardsIOAPIMock{}
	da := RailDeviceAPI{boardsIOAPI: ba}
	da.devices = make(map[string]struct{})
	// act
	_, err := da.createTurnout(devicerecipe.Ingredients{BoardID: "error", BoardPinNrSec: 88})
	// assert
	require.NotNil(err)
	assert.Equal("test error", err.Error())
}

func (am boardsIOAPIMock) GetInputPin(boardID string, boardPinNr uint8) (boardPin *boardpin.Input, err error) {
	if boardID == "error" {
		err = fmt.Errorf("test error")
	}
	boardPin = &boardpin.Input{}
	return
}

func (am boardsIOAPIMock) GetOutputPin(boardID string, boardPinNr uint8) (boardPin *boardpin.Output, err error) {
	if boardID == "error" && boardPinNr == 88 {
		err = fmt.Errorf("test error")
	}
	boardPin = &boardpin.Output{}
	return
}

func (i inputerMock) RailDeviceName() string                                   { return "test_input" }
func (i inputerMock) StateChanged(visitor string) (hasChanged bool, err error) { return }
func (i inputerMock) IsOn() bool                                               { return false }

func (r runnerMock) RailDeviceName() string                                   { return r.name }
func (r runnerMock) SwitchOn() (err error)                                    { return }
func (r runnerMock) SwitchOff() (err error)                                   { return }
func (r runnerMock) StateChanged(visitor string) (hasChanged bool, err error) { return }
func (r runnerMock) IsOn() bool                                               { return false }
