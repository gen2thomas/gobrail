package devicerecipe

import (
	"fmt"
	"time"
)

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
	// TypUnknown is fo fallback
	TypUnknown
)

// TypeMap is the string representation to the underlying int type
var TypeMap = map[string]railDeviceType{
	"Button": Button, "ToggleButton": ToggleButton,
	"Lamp": Lamp, "TwoLightsSignal": TwoLightsSignal, "Turnout": Turnout,
	"TypUnknown": TypUnknown,
}

// RailDeviceRecipe describes a recipe to creat an new rail device
type RailDeviceRecipe struct {
	Name           string `json:"Name"`
	Type           string `json:"Type"`
	BoardID        string `json:"BoardID"`
	BoardPinNrPrim uint8  `json:"BoardPinNrPrim"`
	BoardPinNrSec  uint8  `json:"BoardPinNrSec"`
	StartingDelay  string `json:"StartingDelay"`
	StoppingDelay  string `json:"StoppingDelay"`
	Connect        string `json:"Connect"`
}

// RailPlan represents all recipes for rail devices
type RailPlan struct {
	DeviceRecipes []RailDeviceRecipe `json:"DeviceRecipes"`
}

// Verify is checking the parsability of string values to the corresponding type
func (r RailDeviceRecipe) Verify() (err error) {
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
func (r *RailDeviceRecipe) FillEmptyDefaults() {
	if r.StartingDelay == "" {
		r.StartingDelay = "0"
	}
	if r.StoppingDelay == "" {
		r.StoppingDelay = "0"
	}
	return
}

func (r RailDeviceRecipe) String() string {
	return fmt.Sprintf("Name: %s, Type: %s, BoardID: %s, BoardPinNr: %d, BoardPinNrSecond: %d, StartingDelay: %s, StoppingDelay: %s, Connect: %s",
		r.Name, r.Type, r.BoardID, r.BoardPinNrPrim, r.BoardPinNrSec, r.StartingDelay, r.StoppingDelay, r.Connect)
}
