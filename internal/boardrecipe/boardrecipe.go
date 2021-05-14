package boardrecipe

// A boardrecipe is the description how to create an board

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gen2thomas/gobrail/internal/errwrap"
	"github.com/gen2thomas/gobrail/internal/jsonrecipe"
)

// Schema is for json validation
var Schema = "./schemas/board.schema.json"

type boardType uint8

const (
	// Type2i is the board with a single PCA9501 with 8 inputs
	Type2i boardType = iota
	// Type2o is the board with a single PCA9501 with 8 amplified outputs
	Type2o boardType = iota
	// Type2io is the board with a single PCA9501 with 4 inputs and 4 amplified outputs
	Type2io boardType = iota
	// TypUnknown is for fall back
	TypUnknown
)

// TypeMap is the string representation to the underlying "boardType"
var TypeMap = map[string]boardType{
	"Type2i": Type2i, "Type2o": Type2o, "Type2io": Type2io, "TypUnknown": TypUnknown,
}

// Ingredients is a short description to create a new board
type Ingredients struct {
	Name        string `json:"Name"`
	Type        string `json:"Type"`
	ChipDevAddr uint8  `json:"ChipDevAddr"`
}

// ReadIngredients is parsing json board description to a board recipe
func ReadIngredients(boardFile string) (recipe Ingredients, err error) {
	boardFile, err = jsonrecipe.PrepareAndValidate(Schema, boardFile)
	if err != nil {
		return
	}

	var jsonFile *os.File
	var byteValue []byte
	jsonFile, err = os.Open(boardFile)
	if err == nil {
		byteValue, err = ioutil.ReadAll(jsonFile)
	}
	if err == nil {
		err = json.Unmarshal(byteValue, &recipe)
	}
	err = errwrap.Wrap(err, jsonFile.Close())
	if err == nil {
		err = recipe.verify()
	}
	if err != nil {
		err = fmt.Errorf("%s for file %s", err.Error(), boardFile)
	}
	return
}

// Verify is checking that string values are parsable to the corresponding type
func (r Ingredients) verify() (err error) {
	// check for type string is known
	if _, ok := TypeMap[r.Type]; !ok {
		err = fmt.Errorf("The given type '%s' is unknown", r.Type)
	}
	return
}

func (r Ingredients) String() string {
	return fmt.Sprintf("Name: %s, Type: %s, Chip address: %d", r.Name, r.Type, r.ChipDevAddr)
}
