package boardrecipe

// A boardrecipe is the description how to create an board

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type boardType uint8

const (
	// Typ2 is the board with a single PCA9501 and 4 amplified outputs
	Typ2 boardType = iota
	// TypUnknown is fo fallback
	TypUnknown
)

// TypeMap is the string representation to the underlying "boardType"
var TypeMap = map[string]boardType{
	"Typ2": Typ2, "TypUnknown": TypUnknown,
}

// Ingredients is a short description to create a new board
type Ingredients struct {
	Name        string `json:"Name"`
	Type        string `json:"Type"`
	ChipDevAddr uint8  `json:"ChipDevAddr"`
}

// ReadIngredients is parsing json board description to a board recipe
func ReadIngredients(boardFile string) (recipe Ingredients, err error) {
	boardFile, err = filepath.Abs(boardFile)
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

// Verify is checking the parsability of string values to the corresponding type
func (r Ingredients) Verify() (err error) {
	// check for type string is known
	if _, ok := TypeMap[r.Type]; !ok {
		err = fmt.Errorf("The given type '%s' is unknown", r.Type)
	}
	return
}
