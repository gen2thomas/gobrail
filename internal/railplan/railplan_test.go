package railplan

import (
	"path/filepath"
	"testing"

	"github.com/gen2thomas/gobrail/internal/boardrecipe"
	"github.com/gen2thomas/gobrail/internal/devicerecipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const recipesBase = "../../test/data/"

func TestReadCookBook(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	cookbook := recipesBase + "plans/plan_2boards_2devices_test.json"
	//cookbook := recipesBase + "plans/plan.json"
	oldSchema := schema
	schema, _ = filepath.Abs("../../schemas/plan.schema.json")
	defer func() { schema = oldSchema }()
	// act
	book, err := ReadCookBook(cookbook)
	// assert
	require.Nil(err)
	require.NotNil(book)
	assert.Equal(2, len(book.DeviceRecipes))
	assert.Equal(2, len(book.BoardRecipes))
}

func TestAddBoardRecipe(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	oldSchema := schema
	boardrecipe.Schema, _ = filepath.Abs("../../schemas/board.schema.json")
	defer func() { boardrecipe.Schema = oldSchema }()
	recipe := recipesBase + "boardrecipes/board_test.json"
	book := &CookBook{}
	// act
	err := book.AddBoardRecipe(recipe)
	// assert
	require.Nil(err)
	require.NotNil(book)
	assert.Equal(0, len(book.DeviceRecipes))
	require.Equal(1, len(book.BoardRecipes))
	assert.Equal("B1", book.BoardRecipes[0].Name)
	// other stuff is tested by "boardrecipe_test.go"
}

func TestAddDeviceRecipe(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	oldSchema := schema
	devicerecipe.Schema, _ = filepath.Abs("../../schemas/raildevice.schema.json")
	defer func() { devicerecipe.Schema = oldSchema }()
	recipe := recipesBase + "devicerecipes/device_test.json"
	book := &CookBook{}
	// act
	err := book.AddDeviceRecipe(recipe)
	// assert
	require.Nil(err)
	require.NotNil(book)
	assert.Equal(0, len(book.BoardRecipes))
	require.Equal(1, len(book.DeviceRecipes))
	assert.Equal("D1", book.DeviceRecipes[0].Name)
	// other stuff is tested by "devicerecipe_test.go"
}
