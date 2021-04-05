package boardsapi_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gobot.io/x/gobot/drivers/i2c"

	"github.com/gen2thomas/gobrail/internal/boardsapi"
)

type adaptorMock struct {
	name string
}

var boardRecipeTyp2 = boardsapi.BoardRecipe{
	Name:        "TestRecipeTyp2",
	ChipDevAddr: 0x07,
	BoardType:   boardsapi.Typ2,
}

var boardRecipeUnknown = boardsapi.BoardRecipe{
	Name:        "TestRecipeTypUnknown",
	ChipDevAddr: 0x27,
	BoardType:   boardsapi.TypUnknown,
}

func TestNewBoardsAPI(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	// act
	api := boardsapi.NewBoardsAPI(new(adaptorMock), []boardsapi.BoardRecipe{boardRecipeTyp2})
	// assert
	require.NotNil(api)
	assert.Equal("Boards: 1\nBoard Id: TestRecipeTyp2, Name: TestRecipeTyp2, Chips: 1, Pins: 16\n\nNo pins mapped yet", fmt.Sprintf("%s", api))
}

func TestNewBoardsAPIWithUnknownTypeGetsEmptyBoards(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	// act
	api := boardsapi.NewBoardsAPI(new(adaptorMock), []boardsapi.BoardRecipe{boardRecipeUnknown})
	// assert
	require.NotNil(api)
	assert.Equal("No Boards\nNo pins mapped yet", fmt.Sprintf("%s", api))
}

func (a *adaptorMock) GetConnection(address int, bus int) (device i2c.Connection, err error) { return }
func (a *adaptorMock) GetDefaultBus() int                                                    { return 0 }
