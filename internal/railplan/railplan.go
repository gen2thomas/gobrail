package railplan

// A railplan is the description how to create a model railroad controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gen2thomas/gobrail/internal/boardrecipe"
	"github.com/gen2thomas/gobrail/internal/devicerecipe"
	"github.com/gen2thomas/gobrail/internal/errwrap"
	"github.com/gen2thomas/gobrail/internal/jsonrecipe"
)

// TODO: can write json plan from plan-object-list of creator

const schema = "./schemas/plan.schema.json"

// CookBook contains all recipes for boards and rail devices
type CookBook struct {
	DeviceRecipes []devicerecipe.Ingredients `json:"DeviceRecipes"`
	BoardRecipes  []boardrecipe.Ingredients  `json:"BoardRecipes"`
}

// ReadCookBook is parsing json plan to a list of device recipes
func ReadCookBook(planFile string) (railPlan CookBook, err error) {
	planFile, err = jsonrecipe.PrepareAndValidate(schema, planFile)
	if err != nil {
		return
	}

	var jsonFile *os.File
	var byteValue []byte
	jsonFile, err = os.Open(planFile)
	if err == nil {
		byteValue, err = ioutil.ReadAll(jsonFile)
	}
	if err == nil {
		err = json.Unmarshal(byteValue, &railPlan)
	}
	err = errwrap.Wrap(err, jsonFile.Close())
	if err != nil {
		err = fmt.Errorf("%s for file %s", err.Error(), planFile)
	}

	return
}

// AddBoardRecipe read and add a board to menu card
func (p *CookBook) AddBoardRecipe(boardFile string) (err error) {
	var recipe boardrecipe.Ingredients
	if recipe, err = boardrecipe.ReadIngredients(boardFile); err != nil {
		return
	}
	p.BoardRecipes = append(p.BoardRecipes, recipe)
	return
}

// AddDeviceRecipe read and add a rail device to menu card
func (p *CookBook) AddDeviceRecipe(deviceFile string) (err error) {
	var recipe devicerecipe.Ingredients
	if recipe, err = devicerecipe.ReadIngredients(deviceFile); err != nil {
		return
	}
	p.DeviceRecipes = append(p.DeviceRecipes, recipe)
	return
}
