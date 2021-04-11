package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gobot.io/x/gobot"

	"github.com/gen2thomas/gobrail/internal/boardpin"
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
		0:  {PinType: boardpin.Binary},
		1:  {PinType: boardpin.Memory},
		2:  {PinType: boardpin.Analog},
		3:  {PinType: boardpin.NBinary},
		4:  {PinType: boardpin.NBinaryR},
		5:  {PinType: boardpin.NBinaryW},
		6:  {PinType: boardpin.MemoryR},
		7:  {PinType: boardpin.AnalogR},
		8:  {PinType: boardpin.MemoryW},
		9:  {PinType: boardpin.MemoryR},
		10: {PinType: boardpin.BinaryR},
		11: {PinType: boardpin.AnalogW},
		12: {PinType: boardpin.BinaryW},
	}

	testBoard := &Board{
		pins: boardPins,
	}

	// act
	pinsAll := testBoard.GetPinNumbers()
	pinsBin := testBoard.GetPinNumbersOfType(boardpin.Binary, boardpin.BinaryW, boardpin.BinaryR, boardpin.NBinary, boardpin.NBinaryW, boardpin.NBinaryR)
	pinsAna := testBoard.GetPinNumbersOfType(boardpin.Analog, boardpin.AnalogW, boardpin.AnalogR)
	pinsMem := testBoard.GetPinNumbersOfType(boardpin.Memory, boardpin.MemoryW, boardpin.MemoryR)

	// assert
	assert.Equal(6, len(pinsBin))
	assert.Equal(3, len(pinsAna))
	assert.Equal(4, len(pinsMem))
	assert.Equal(13, len(pinsAll))
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

func (d *deviceMock) Name() string                                                      { return d.name }
func (d *deviceMock) SetName(s string)                                                  { d.name = s }
func (d *deviceMock) Start() (err error)                                                { return }
func (d *deviceMock) Halt() (err error)                                                 { return }
func (d *deviceMock) Connection() gobot.Connection                                      { return nil }
func (d *deviceMock) WriteGPIO(pin uint8, val uint8) (err error)                        { return }
func (d *deviceMock) ReadGPIO(pin uint8) (val uint8, err error)                         { return }
func (d *deviceMock) Command(string) (command func(map[string]interface{}) interface{}) { return }
