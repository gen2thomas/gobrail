package devicerecipe

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

type fillTest struct {
	di   Ingredients
	want Ingredients
}

func Test_verify(t *testing.T) {
	var verifyTests = map[string]verifyTest{
		"WrongType":       {di: Ingredients{Type: "WrongType"}, wantErr: "type 'WrongType' is unknown"},
		"WrongStartDelay": {di: Ingredients{Type: "Button", StartingDelay: "WrongStartDelay"}, wantErr: "start delay 'WrongStartDelay' is not parsable"},
		"WrongStopDelay":  {di: Ingredients{Type: "Button", StartingDelay: "1m", StoppingDelay: "WrongStopDelay"}, wantErr: "stop delay 'WrongStopDelay' is not parsable"},
		"NoError":         {di: Ingredients{Type: "Button", StartingDelay: "1m", StoppingDelay: "1s"}},
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

func Test_fillEmptyDefaults(t *testing.T) {
	var fillTests = map[string]fillTest{
		"EmptyStart": {di: Ingredients{StartingDelay: "", StoppingDelay: "1m"}, want: Ingredients{StartingDelay: "0", StoppingDelay: "1m"}},
		"EmptyStop":  {di: Ingredients{StartingDelay: "2h", StoppingDelay: ""}, want: Ingredients{StartingDelay: "2h", StoppingDelay: "0"}},
		"NoFill":     {di: Ingredients{StartingDelay: "2h1m", StoppingDelay: "3m2s"}, want: Ingredients{StartingDelay: "2h1m", StoppingDelay: "3m2s"}},
	}
	for name, ft := range fillTests {
		t.Run(name, func(t *testing.T) {
			// arrange
			assert := assert.New(t)
			// act
			ft.di.fillEmptyDefaults()
			// assert
			assert.Equal(ft.want.StartingDelay, ft.di.StartingDelay)
			assert.Equal(ft.want.StoppingDelay, ft.di.StoppingDelay)
		})
	}
}

func TestReadIngredients(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	oldSchema := Schema
	Schema, _ = filepath.Abs("../../schemas/raildevice.schema.json")
	defer func() { Schema = oldSchema }()
	recipe := recipesBase + "devicerecipes/device_test.json"
	// act
	ing, err := ReadIngredients(recipe)
	// assert
	require.Nil(err)
	require.NotNil(ing)
	assert.Equal("D1", ing.Name)
	assert.Equal("Turnout", ing.Type)
	assert.Equal("B1", ing.BoardID)
	assert.Equal(uint8(1), ing.BoardPinNrPrim)
	assert.Equal(uint8(2), ing.BoardPinNrSec)
	assert.Equal("0.1s", ing.StartingDelay)
	assert.Equal("0.15s", ing.StoppingDelay)
	assert.Equal("D2", ing.Connect)
}
