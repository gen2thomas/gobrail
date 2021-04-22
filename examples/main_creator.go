// +build example
//
// Do not build by default.

package main

import (	
	"github.com/gen2thomas/gobrail/internal/app/gobrailcreator"
)

// For a breadboard schematic refer to docs/images/PCA9501_Lamps_Buttons.png
// Just substidude the magnets with LED's and a 150Ohm resistor.

func main() {
	adaptype, _ := gobrailcreator.ParseAdaptorType("digispark")
	recipes := gobrailcreator.RecipeFiles{
		boards: []string{"./test/data/board_typ2_0x04.json", "./test/data/board_typ2_0x05.json"}
		devices: []string{"./test/data/device_button4.json", "./test/data/device_togglebutton5.json"}
	}
	gobrailcreator.Create(false, "dummy rob name", adaptype, "./test/data/plan1.json", recipes)
}
