package boardrecipe

// A boardrecipe is the description how to create an board

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gen2thomas/gobrail/internal/jsonrecipe"
)

const schema = "./schemas/board.schema.json"

type boardType uint8

const (
	// Type2 is the board with a single PCA9501 and 4 amplified outputs
	Type2 boardType = iota
	// TypUnknown is for fall back
	TypUnknown
)

// TypeMap is the string representation to the underlying "boardType"
var TypeMap = map[string]boardType{
	"Type2": Type2, "TypUnknown": TypUnknown,
}

// Ingredients is a short description to create a new board
type Ingredients struct {
	Name        string `json:"Name"`
	Type        string `json:"Type"`
	ChipDevAddr uint8  `json:"ChipDevAddr"`
}

// ReadIngredients is parsing json board description to a board recipe
func ReadIngredients(boardFile string) (recipe Ingredients, err error) {
	boardFile, err = jsonrecipe.PrepareAndValidate(schema, boardFile)
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
	if err2 := jsonFile.Close(); err2 != nil {
		if err == nil {
			err = err2
			return
		}
		err = fmt.Errorf("%s for file %s %w", err.Error(), boardFile, err2)
	}
	return
}

// Verify is checking that string values are parsable to the corresponding type
func (r Ingredients) Verify() (err error) {
	// check for type string is known
	if _, ok := TypeMap[r.Type]; !ok {
		err = fmt.Errorf("The given type '%s' is unknown", r.Type)
	}
	return
}

func (r Ingredients) String() string {
	return fmt.Sprintf("Name: %s, Type: %s, Chip address: %d", r.Name, r.Type, r.ChipDevAddr)
}
