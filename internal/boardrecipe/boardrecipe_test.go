package boardrecipe

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type verifyTest struct {
	di      Ingredients
	wantErr string
}

func Test_verify(t *testing.T) {
	var verifyTests = map[string]verifyTest{
		"WrongType": {di: Ingredients{Type: "WrongType"}, wantErr: "type 'WrongType' is unknown"},
		"NoError":   {di: Ingredients{Type: "Type2"}},
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
