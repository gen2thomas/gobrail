package devicerecipe

// A devicerecipe is the description how to create an rail device

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gen2thomas/gobrail/internal/jsonrecipe"
)

const schema = "./schemas/raildevice.schema.json"

type railDeviceType uint8

const (
	// Button is a input device with one input
	Button railDeviceType = iota
	// ToggleButton is a input device with one input
	ToggleButton
	// Lamp is a output device with one output
	Lamp
	// TwoLightsSignal is a output device with two outputs, both outputs can't have the same state
	TwoLightsSignal
	// Turnout is a output device with two outputs
	Turnout
	// TypUnknown is for fall-back
	TypUnknown
)

// TypeMap is the string representation to the underlying "railDeviceType"
var TypeMap = map[string]railDeviceType{
	"Button": Button, "ToggleButton": ToggleButton,
	"Lamp": Lamp, "TwoLightsSignal": TwoLightsSignal, "Turnout": Turnout,
	"TypUnknown": TypUnknown,
}

// Ingredients describes a recipe to create an new rail device
type Ingredients struct {
	Name           string `json:"Name"`
	Type           string `json:"Type"`
	BoardID        string `json:"BoardID"`
	BoardPinNrPrim uint8  `json:"BoardPinNrPrim"`
	BoardPinNrSec  uint8  `json:"BoardPinNrSec"`
	StartingDelay  string `json:"StartingDelay"`
	StoppingDelay  string `json:"StoppingDelay"`
	Connect        string `json:"Connect"`
}

// TODO: json verification
// TODO: wrapped errors
// TODO: can write json single object description from a a plan-object

// ReadIngredients is parsing json device description to a device recipe
func ReadIngredients(deviceFile string) (recipe Ingredients, err error) {
	deviceFile, err = jsonrecipe.PrepareAndValidate(schema, deviceFile)
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

// Verify is checking that string values are parsable to the corresponding type
func (r Ingredients) Verify() (err error) {
	// check for type string is known
	if _, ok := TypeMap[r.Type]; !ok {
		err = fmt.Errorf("The given type '%s' is unknown", r.Type)
	}
	// check for delays are parsable
	if _, err1 := time.ParseDuration(r.StartingDelay); err1 != nil {
		err = fmt.Errorf("The given start delay '%s' is not parsable, %w", r.StartingDelay, err)
	}
	if _, err1 := time.ParseDuration(r.StoppingDelay); err1 != nil {
		err = fmt.Errorf("The given stop delay '%s' is not parsable, %w", r.StoppingDelay, err)
	}

	return
}

// FillEmptyDefaults will correct some optional values after parsing
func (r *Ingredients) FillEmptyDefaults() {
	if r.StartingDelay == "" {
		r.StartingDelay = "0"
	}
	if r.StoppingDelay == "" {
		r.StoppingDelay = "0"
	}
	return
}

func (r Ingredients) String() string {
	return fmt.Sprintf("Name: %s, Type: %s, BoardID: %s, BoardPinNr: %d, BoardPinNrSecond: %d, StartingDelay: %s, StoppingDelay: %s, Connect: %s",
		r.Name, r.Type, r.BoardID, r.BoardPinNrPrim, r.BoardPinNrSec, r.StartingDelay, r.StoppingDelay, r.Connect)
}
