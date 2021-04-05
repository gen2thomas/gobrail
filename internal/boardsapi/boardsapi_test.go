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
	require := require.New(t)
	// act
	api := NewBoardsAPI(new(adaptorMock), []BoardRecipe{boardRecipeTyp2})
	// assert
	require.NotNil(api)
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

func TestGetFreeAPIPinsWithoutBoardGetsEmptyList(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	api := &BoardsAPI{
		mappedPins: make(APIPinsMap),
		boards:     make(BoardsMap),
	}
	// act
	fp := api.GetFreeAPIPins("NoExistend", board.Binary)
	// assert
	require.NotNil(fp)
	assert.Equal(0, len(fp))
}

func TestGetFreeAPIPins(t *testing.T) {
	// arrange
	assert := assert.New(t)
	b1 := board.NewBoardForTestWithoutChips("TestBoard1", 2, 5, 1)
	bm := make(BoardsMap)
	bm["TestBoard1"] = b1
	api := &BoardsAPI{
		mappedPins: make(APIPinsMap),
		boards:     bm,
	}
	// act
	fp1bin := api.GetFreeAPIPins("TestBoard1", board.Binary)
	fp1ana := api.GetFreeAPIPins("TestBoard1", board.Analog)
	fp1mem := api.GetFreeAPIPins("TestBoard1", board.Memory)
	// assert
	assert.Equal(2, len(fp1bin))
	assert.Equal(5, len(fp1ana))
	assert.Equal(1, len(fp1mem))
}

func TestGetMappedAPIPinsWithoutBoardGetsEmptyList(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	api := &BoardsAPI{
		mappedPins: make(APIPinsMap),
		boards:     make(BoardsMap),
	}
	// act
	mp := api.GetMappedAPIPins("NoExistend", board.Binary)
	// assert
	require.NotNil(mp)
	assert.Equal(0, len(mp))
}

func TestGetMappedAPIPins(t *testing.T) {
	// arrange
	assert := assert.New(t)
	b1 := board.NewBoardForTestWithoutChips("TestBoard1", 3, 3, 3)
	bm := make(BoardsMap)
	bm["TestBoard1"] = b1
	mps := make(APIPinsMap)
	mps["k1"] = &apiPin{boardID: "TestBoard1", boardPinNr: 4, railDeviceName: "ana device 1"}
	mps["k2"] = &apiPin{boardID: "TestBoard1", boardPinNr: 1, railDeviceName: "binary device 1"}
	mps["k3"] = &apiPin{boardID: "TestBoard1", boardPinNr: 5, railDeviceName: "ana device 2"}
	mps["k4"] = &apiPin{boardID: "TestBoard1", boardPinNr: 8, railDeviceName: "mem device"}
	mps["k5"] = &apiPin{boardID: "TestBoard1", boardPinNr: 0, railDeviceName: "binary device 2"}
	mps["k6"] = &apiPin{boardID: "TestBoard1", boardPinNr: 2, railDeviceName: "binary device 3"}
	api := &BoardsAPI{
		mappedPins: mps,
		boards:     bm,
	}
	// act
	mp1bin := api.GetMappedAPIPins("TestBoard1", board.Binary)
	mp1ana := api.GetMappedAPIPins("TestBoard1", board.Analog)
	mp1mem := api.GetMappedAPIPins("TestBoard1", board.Memory)
	// assert
	assert.Equal(3, len(mp1bin))
	assert.Equal(2, len(mp1ana))
	assert.Equal(1, len(mp1mem))
}

func (a *adaptorMock) GetConnection(address int, bus int) (device i2c.Connection, err error) { return }
func (a *adaptorMock) GetDefaultBus() int                                                    { return 0 }
