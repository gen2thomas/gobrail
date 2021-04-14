package devicerecipe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// TODO: json verification
// TODO: wrapped errors

// ReadPlan is parsing json plan to a list of device recipes
func ReadPlan(planFile string) (recipes []RailDeviceRecipe, err error) {
	planFile, err = filepath.Abs(planFile)
	if err != nil {
		return
	}

	var jsonFile *os.File
	var byteValue []byte
	jsonFile, err = os.Open(planFile)
	if err == nil {
		byteValue, err = ioutil.ReadAll(jsonFile)
	}
	var railPlan RailPlan
	if err == nil {
		err = json.Unmarshal(byteValue, &railPlan)
	}
	err2 := jsonFile.Close()
	if err == nil {
		for _, recipe := range railPlan.DeviceRecipes {
			recipe.FillEmptyDefaults()
			recipes = append(recipes, recipe)
		}
		err = err2
		return
	}
	err = fmt.Errorf("%s for file %s %w", err.Error(), planFile, err2)
	return
}

// ReadDevice is parsing json device description to a device recipe
func ReadDevice(deviceFile string) (recipe RailDeviceRecipe, err error) {
	deviceFile, err = filepath.Abs(deviceFile)
	if err != nil {
		return
	}

	var jsonFile *os.File
	var byteValue []byte
	jsonFile, err = os.Open(deviceFile)
	if err == nil {
		byteValue, err = ioutil.ReadAll(jsonFile)
	}
	if err == nil {
		err = json.Unmarshal(byteValue, &recipe)
	}
	err2 := jsonFile.Close()
	if err == nil {
		recipe.FillEmptyDefaults()
		err = err2
		return
	}
	err = fmt.Errorf("%s for file %s %w", err.Error(), deviceFile, err2)
	return
}
