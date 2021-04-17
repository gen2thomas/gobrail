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

type rwTest struct {
	pType  boardpin.PinType
	fails  bool
	expVal uint8
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

func TestWriteValue(t *testing.T) {
	// arrange
	assert := assert.New(t)
	// note: analog is not implemented yet, therfore fails
	var wTests = []rwTest{
		{pType: boardpin.Binary, fails: false, expVal: uint8(1)},
		{pType: boardpin.BinaryR, fails: true, expVal: uint8(1)},
		{pType: boardpin.BinaryW, fails: false, expVal: uint8(0)},
		{pType: boardpin.NBinary, fails: false, expVal: uint8(0)},
		{pType: boardpin.NBinaryR, fails: true, expVal: uint8(0)},
		{pType: boardpin.NBinaryW, fails: false, expVal: uint8(0)},
		{pType: boardpin.Memory, fails: false, expVal: uint8(1)},
		{pType: boardpin.MemoryR, fails: true, expVal: uint8(1)},
		{pType: boardpin.MemoryW, fails: false, expVal: uint8(0)},
		{pType: boardpin.Analog, fails: true, expVal: uint8(0)},
		{pType: boardpin.AnalogR, fails: true, expVal: uint8(0)},
		{pType: boardpin.AnalogW, fails: true, expVal: uint8(0)},
	}
	for _, wt := range wTests {
		name := "for " + boardpin.PinTypeMsgMap[wt.pType]
		t.Run(name, func(t *testing.T) {
			// arrange
			cID := "chipNR"
			pNr := uint8(3)
			d := &deviceMock{name: "dev1"}
			b := &Board{}
			b.chips = map[string]*chip{cID: {driver: d}}
			b.pins = PinsMap{pNr: {ChipID: cID, PinType: wt.pType}}
			// act
			err := b.WriteValue(pNr, wt.expVal)
			// assert
			assert.Equal(wt.fails, err != nil)
		})
	}
}

func TestReadValue(t *testing.T) {
	// arrange
	assert := assert.New(t)
	// note: analog is not implemented yet, therfore fails
	var rTests = []rwTest{
		{pType: boardpin.Binary, fails: false, expVal: uint8(1)},
		{pType: boardpin.BinaryR, fails: false, expVal: uint8(1)},
		{pType: boardpin.BinaryW, fails: true, expVal: uint8(0)},
		{pType: boardpin.NBinary, fails: false, expVal: uint8(0)},
		{pType: boardpin.NBinaryR, fails: false, expVal: uint8(0)},
		{pType: boardpin.NBinaryW, fails: true, expVal: uint8(0)},
		{pType: boardpin.Memory, fails: false, expVal: uint8(1)},
		{pType: boardpin.MemoryR, fails: false, expVal: uint8(1)},
		{pType: boardpin.MemoryW, fails: true, expVal: uint8(0)},
		{pType: boardpin.Analog, fails: true, expVal: uint8(0)},
		{pType: boardpin.AnalogR, fails: true, expVal: uint8(0)},
		{pType: boardpin.AnalogW, fails: true, expVal: uint8(0)},
	}
	for _, rt := range rTests {
		name := "for " + boardpin.PinTypeMsgMap[rt.pType]
		t.Run(name, func(t *testing.T) {
			// arrange
			cID := "chipNR"
			pNr := uint8(3)
			d := &deviceMock{name: "dev1"}
			b := &Board{}
			b.chips = map[string]*chip{cID: {driver: d}}
			b.pins = PinsMap{pNr: {ChipID: cID, PinType: rt.pType}}
			// act
			val, err := b.ReadValue(pNr)
			// assert
			assert.Equal(rt.fails, err != nil)
			assert.Equal(rt.expVal, val)
		})
	}
}

func (d *deviceMock) Name() string                               { return d.name }
func (d *deviceMock) SetName(s string)                           { d.name = s }
func (d *deviceMock) Start() (err error)                         { return }
func (d *deviceMock) Halt() (err error)                          { return }
func (d *deviceMock) Connection() gobot.Connection               { return nil }
func (d *deviceMock) WriteGPIO(pin uint8, val uint8) (err error) { return }
func (d *deviceMock) ReadGPIO(pin uint8) (val uint8, err error)  { return }
func (d *deviceMock) Command(string) (command func(map[string]interface{}) interface{}) {
	command = func(map[string]interface{}) interface{} {
		return map[string]interface{}{"err": nil, "val": uint8(1)}
	}
	return
}
