package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot"
)

type deviceMock struct {
	name       string
	connection gobot.Connection
}

type adaptorMock struct {
	name string
}

func TestDevices(t *testing.T) {
	// arrange
	assert := assert.New(t)
	dev1 := &deviceMock{name: "dev1"}
	dev2 := &deviceMock{name: "dev2"}

	testBoard := &Board{
		chips: map[string]*chip{
			"testchip1": {
				device: dev1,
			},
			"testchip2": {
				device: dev2,
			},
		},
	}

	// act
	devs := testBoard.Devices()

	// assert
	assert.Equal(len(devs), 2)
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
	assert.Equal(len(pinsBin), 3)
	assert.Equal(len(pinsAna), 2)
	assert.Equal(len(pinsMem), 1)
}

func (d *deviceMock) Name() string                                                      { return d.name }
func (d *deviceMock) SetName(s string)                                                  { d.name = s }
func (d *deviceMock) Start() (err error)                                                { return }
func (d *deviceMock) Halt() (err error)                                                 { return }
func (d *deviceMock) Connection() gobot.Connection                                      { return d.connection }
func (d *deviceMock) WriteGPIO(pin uint8, val uint8) (err error)                        { return }
func (d *deviceMock) ReadGPIO(pin uint8) (val uint8, err error)                         { return }
func (d *deviceMock) Command(string) (command func(map[string]interface{}) interface{}) { return }

func (a *adaptorMock) Name() string          { return a.name }
func (a *adaptorMock) SetName(n string)      { a.name = n }
func (a *adaptorMock) Connect() (err error)  { return }
func (a *adaptorMock) Finalize() (err error) { return }
