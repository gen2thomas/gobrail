package boardsapi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gobot.io/x/gobot/drivers/i2c"

	"github.com/gen2thomas/gobrail/internal/board"
)

type adaptorMock struct {
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
	b1 := board.NewBoardForTestWithoutChips("TestBoard1", 1, 1, 1)
	bm := make(BoardsMap)
	bm["TestBoard1"] = b1
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
	b1 := board.NewBoardForTestWithoutChips("TestBoard1", 1, 1, 1)
	bm := make(BoardsMap)
	bm["TestBoard1"] = b1
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
