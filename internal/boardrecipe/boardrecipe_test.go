package boardrecipe

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const recipesBase = "../../test/data/"

type verifyTest struct {
	di      Ingredients
	wantErr string
}

func Test_verify(t *testing.T) {
	var verifyTests = map[string]verifyTest{
		"WrongType": {di: Ingredients{Type: "WrongType"}, wantErr: "type 'WrongType' is unknown"},
		"NoError":   {di: Ingredients{Type: "Type2io"}},
	}
	for name, vt := range verifyTests {
		t.Run(name, func(t *testing.T) {
			// arrange
			assert := assert.New(t)
			require := require.New(t)
			// act
			err := vt.di.verify()
			// assert
			if vt.wantErr == "" {
				assert.Nil(err)
			} else {
				require.NotNil(err)
				assert.Contains(err.Error(), vt.wantErr)
			}
		})
	}
}

func TestReadIngredients(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	oldSchema := Schema
	Schema, _ = filepath.Abs("../../schemas/board.schema.json")
	defer func() { Schema = oldSchema }()
	recipe := recipesBase + "boardrecipes/board_test.json"
	// act
	ing, err := ReadIngredients(recipe)
	// assert
	require.Nil(err)
	require.NotNil(ing)
	assert.Equal("B1", ing.Name)
	assert.Equal("Type2io", ing.Type)
	assert.Equal(uint8(1), ing.ChipDevAddr)
}
