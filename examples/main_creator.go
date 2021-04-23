// +build example
//
// Do not build by default.

package main

import (
	"fmt"

	"github.com/gen2thomas/gobrail/internal/app/gobrailcreator"
)

// For a breadboard schematic refer to docs/images/PCA9501_Lamps_Buttons.png
// Just substitute the magnets with LED's and a 150Ohm resistor.

func main() {
	adaptype, _ := gobrailcreator.ParseAdaptorType("digispark")
	recipes := gobrailcreator.RecipeFiles{
		Boards: []string{"./test/data/board_typ2_0x04.json", "./test/data/board_typ2_0x05.json"},
		Devices: []string{"./test/data/device_button4.json", "./test/data/device_togglebutton5.json"},
	}
	if _, err := gobrailcreator.Create(false, "dummy rob name", adaptype, "./test/data/plan1.json", recipes); err != nil{
		fmt.Println("Error occurred:", err)
	}
}
