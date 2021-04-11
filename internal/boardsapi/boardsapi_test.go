package boardsapi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"

	"github.com/gen2thomas/gobrail/internal/boardpin"
)

type adaptorMock struct {
	name string
}

type boardsMock struct {
	name    string
	binPins uint8
	anaPins uint8
	memPins uint8
}

var boardRecipeTyp2 = BoardRecipe{
	Name:        "TestRecipeTyp2",
	ChipDevAddr: 0x07,
	BoardType:   Typ2,
}

var boardRecipeUnknown = BoardRecipe{
	Name:        "TestRecipeTypUnknown",
	ChipDevAddr: 0x27,
	BoardType:   TypUnknown,
}

func TestNewBoardsAPI(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	// act
	api := NewBoardsAPI(new(adaptorMock))
	api.AddBoard(boardRecipeTyp2)
	// assert
	require.Equal(1, len(api.boards))
	require.Equal(1, len(api.usedPins))
	assert.Equal(0, len(api.usedPins[boardRecipeTyp2.Name]))
	assert.Equal("Boards: 1\nBoard Id: TestRecipeTyp2, Name: TestRecipeTyp2, Chips: 1, Pins: 16\n\n", fmt.Sprintf("%s", api))
}

func TestNewBoardsAPIWithUnknownTypeGetsEmptyBoards(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	// act
	api := NewBoardsAPI(new(adaptorMock))
	api.AddBoard(boardRecipeUnknown)
	// assert
	require.NotNil(*api)
	assert.Equal("No Boards\n", fmt.Sprintf("%s", api))
}

func TestGetFreePinsWithoutBoardGetsEmptyList(t *testing.T) {
	// arrange
	assert := assert.New(t)
	api := &BoardsAPI{
		usedPins: make(map[string]boardpin.PinNumbers),
		boards:   make(BoardsMap),
	}
	// act
	fp := api.GetFreePins("NoExistend")
	// assert
	assert.Equal(0, len(fp))
}

func TestGetFreePins(t *testing.T) {
	// arrange
	assert := assert.New(t)
	bm := make(BoardsMap)
	bm["TestBoard1"] = &boardsMock{name: "TestBoard1", binPins: 2, anaPins: 5, memPins: 1}
	api := &BoardsAPI{
		usedPins: make(map[string]boardpin.PinNumbers),
		boards:   bm,
	}
	// act
	fp := api.GetFreePins("TestBoard1")
	// assert
	assert.Equal(8, len(fp))
}

func TestGetUsedPinsWithoutBoardGetsEmptyList(t *testing.T) {
	// arrange
	assert := assert.New(t)
	api := &BoardsAPI{
		usedPins: make(map[string]boardpin.PinNumbers),
		boards:   make(BoardsMap),
	}
	// act
	up := api.GetUsedPins("NoExistend")
	// assert
	assert.Equal(0, len(up))
}

func TestGetUsedPins(t *testing.T) {
	// arrange
	assert := assert.New(t)
	api := &BoardsAPI{
		boards:   make(BoardsMap),
		usedPins: make(map[string]boardpin.PinNumbers),
	}
	api.boards["TestBoard1"] = &boardsMock{name: "TestBoard1"}
	api.usedPins["TestBoard1"] = make(boardpin.PinNumbers)
	api.usedPins["TestBoard1"][uint8(1)] = struct{}{}
	api.usedPins["TestBoard1"][uint8(2)] = struct{}{}
	api.usedPins["TestBoard1"][uint8(3)] = struct{}{}
	api.usedPins["TestBoard1"][uint8(4)] = struct{}{}
	api.usedPins["TestBoard1"][uint8(5)] = struct{}{}
	api.usedPins["TestBoard1"][uint8(6)] = struct{}{}
	// act
	usedPins := api.GetUsedPins("TestBoard1")
	// assert
	assert.Equal(6, len(usedPins))
}

func TestGetInputPin(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	api := &BoardsAPI{
		boards:   make(BoardsMap),
		usedPins: make(map[string]boardpin.PinNumbers),
	}
	api.boards["TestBoard"] = &boardsMock{name: "TestBoard", binPins: 1, anaPins: 1, memPins: 1}
	api.usedPins["TestBoard"] = make(boardpin.PinNumbers)
	// act
	pin, err := api.GetInputPin("TestBoard", 2)
	// assert
	require.Nil(err)
	assert.NotNil(pin)
	assert.Equal(1, len(api.usedPins))
}

func TestGetInputPinWhenAlreadyUsedGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	api := &BoardsAPI{
		boards:   make(BoardsMap),
		usedPins: make(map[string]boardpin.PinNumbers),
	}
	api.boards["TestBoard"] = &boardsMock{name: "TestBoard"}
	api.usedPins["TestBoard"] = make(boardpin.PinNumbers)
	api.usedPins["TestBoard"][1] = struct{}{}
	// act
	pin, err := api.GetInputPin("TestBoard", 1)
	// assert
	require.NotNil(err)
	assert.Nil(pin)
	assert.Contains(err.Error(), "already used")
}

func (a *adaptorMock) GetConnection(address int, bus int) (device i2c.Connection, err error) { return }
func (a *adaptorMock) GetDefaultBus() int                                                    { return 0 }

func (b boardsMock) GobotDevices() []gobot.Device { return nil }
func (b boardsMock) GetPinNumbers() boardpin.PinNumbers {
	return createPinNumbersMap(b.binPins + b.anaPins + b.memPins)
}
func (b boardsMock) GetPinNumbersOfType(boardpin.PinType) boardpin.PinNumbers { return nil }
func (b boardsMock) ReadValue(boardPinNr uint8) (uint8, error)                { return 0, nil }
func (b boardsMock) WriteValue(boardPinNr uint8, value uint8) (err error)     { return }
func (b boardsMock) ShowBoardConfig()                                         { return }

func createPinNumbersMap(pinCount uint8) (pinNumbers boardpin.PinNumbers) {
	pinNumbers = make(boardpin.PinNumbers)
	for pinNumber := uint8(0); pinNumber < pinCount; pinNumber++ {
		pinNumbers[pinNumber] = struct{}{}
	}
	return pinNumbers
}
