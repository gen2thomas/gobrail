package boardsapi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
)

type adaptorMock struct {
	name string
}

type boardsMock struct {
	name string
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
	// act
	api := NewBoardsAPI(new(adaptorMock), []BoardRecipe{boardRecipeTyp2})
	// assert
	assert.Equal(1, len(api.boards))
	assert.Equal("Boards: 1\nBoard Id: TestRecipeTyp2, Name: TestRecipeTyp2, Chips: 1, Pins: 16\n\nNo pins mapped yet", fmt.Sprintf("%s", api))
}

func TestNewBoardsAPIWithUnknownTypeGetsEmptyBoards(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	// act
	api := NewBoardsAPI(new(adaptorMock), []BoardRecipe{boardRecipeUnknown})
	// assert
	require.NotNil(*api)
	assert.Equal("No Boards\nNo pins mapped yet", fmt.Sprintf("%s", api))
}

func TestSetValueWhenNotMappedGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	// boards
	bm := make(BoardsMap)
	bm["TestBoard1"] = &boards1Mock{name: "TestBoard1"}
	// mapping
	api := &BoardsAPI{
		mappedPins: make(APIPinsMap),
		boards:     bm,
	}
	// act
	err := api.SetValue("an device", 0)
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "'an device' not mapped yet")
}

func TestReadValueWhenNotMappedGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	// boards
	bm := make(BoardsMap)
	bm["TestBoard1"] = &boards1Mock{name: "TestBoard1"}
	// mapping
	api := &BoardsAPI{
		mappedPins: make(APIPinsMap),
		boards:     bm,
	}
	// act
	_, err := api.GetValue("an device")
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "'an device' not mapped yet")
}

func (a *adaptorMock) GetConnection(address int, bus int) (device i2c.Connection, err error) { return }
func (a *adaptorMock) GetDefaultBus() int                                                    { return 0 }

func (b boardsMock) GobotDevices() []gobot.Device                       { return nil }
func (b boardsMock) GetBinaryPinNumbers() map[uint8]struct{}            { return nil }
func (b boardsMock) GetAnalogPinNumbers() map[uint8]struct{}            { return nil }
func (b boardsMock) GetMemoryPinNumbers() map[uint8]struct{}            { return nil }
func (b boardsMock) ReadValue(boardPinNr uint8) (uint8, error)          { return 0, nil }
func (b boardsMock) SetValue(boardPinNr uint8, value uint8) (err error) { return }
func (b boardsMock) SetAllIoPins() (err error)                          { return }
func (b boardsMock) ResetAllIoPins() (err error)                        { return }
func (b boardsMock) ShowBoardConfig()                                   { return }
