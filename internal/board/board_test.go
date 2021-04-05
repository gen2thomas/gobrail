package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gobot.io/x/gobot"
)

type deviceMock struct {
	name string
}

func TestGobotDevices(t *testing.T) {
	// arrange
	assert := assert.New(t)
	dev1 := &deviceMock{name: "dev1"}
	dev2 := &deviceMock{name: "dev2"}

	testBoard := &Board{
		chips: map[string]*chip{
			"testchip1": {
				driver: dev1,
			},
			"testchip2": {
				driver: dev2,
			},
		},
	}

	// act
	devs := testBoard.GobotDevices()

	// assert
	assert.Equal(2, len(devs))
	assert.Contains(devs, dev1)
	assert.Contains(devs, dev2)
}

func TestPinsOfType(t *testing.T) {
	// arrange
	assert := assert.New(t)
	boardPins := PinsMap{
		0: {pinType: Binary},
		1: {pinType: Memory},
		2: {pinType: Analog},
		3: {pinType: Analog},
		4: {pinType: Binary},
		5: {pinType: Binary},
	}

	testBoard := &Board{
		pins: boardPins,
	}

	// act
	pinsBin := testBoard.PinsOfType(Binary)
	pinsAna := testBoard.PinsOfType(Analog)
	pinsMem := testBoard.PinsOfType(Memory)

	// assert
	assert.Equal(3, len(pinsBin))
	assert.Equal(2, len(pinsAna))
	assert.Equal(1, len(pinsMem))
}

func TestGetBoardPin(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	boardPins := PinsMap{4: {}}
	testBoard := &Board{pins: boardPins}
	// act
	pin, err := testBoard.getBoardPin(4)
	// assert
	require.Nil(err)
	assert.NotNil(pin)
}

func TestGetBoardPinNotThereGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	boardPins := PinsMap{4: {}}
	testBoard := &Board{pins: boardPins}
	// act
	pin, err := testBoard.getBoardPin(3)
	// assert
	assert.NotNil(err)
	assert.Nil(pin)
}

func TestGetDriver(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	dev := deviceMock{name: "Testdriver"}
	bPin := boardPin{chipID: "Testid"}
	boardPins := PinsMap{5: &bPin}
	testBoard := &Board{
		pins:  boardPins,
		chips: map[string]*chip{"Testid": {driver: &dev}},
	}
	// act
	driver, err := testBoard.getDriver(&bPin)
	// assert
	require.Nil(err)
	assert.NotNil(driver)
	assert.Equal("Testdriver", driver.Name())
}

func TestGetDriverNotThereGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	dev := deviceMock{name: "Testdriver"}
	bPin := boardPin{chipID: "Testid1"}
	boardPins := PinsMap{5: &bPin}
	testBoard := &Board{
		pins:  boardPins,
		chips: map[string]*chip{"Testid2": {driver: &dev}},
	}
	// act
	driver, err := testBoard.getDriver(&bPin)
	// assert
	assert.NotNil(err)
	assert.Nil(driver)
}

func (d *deviceMock) Name() string                                                      { return d.name }
func (d *deviceMock) SetName(s string)                                                  { d.name = s }
func (d *deviceMock) Start() (err error)                                                { return }
func (d *deviceMock) Halt() (err error)                                                 { return }
func (d *deviceMock) Connection() gobot.Connection                                      { return nil }
func (d *deviceMock) WriteGPIO(pin uint8, val uint8) (err error)                        { return }
func (d *deviceMock) ReadGPIO(pin uint8) (val uint8, err error)                         { return }
func (d *deviceMock) Command(string) (command func(map[string]interface{}) interface{}) { return }
