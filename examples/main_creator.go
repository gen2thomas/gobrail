// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	
	"github.com/gen2thomas/gobrail/internal/app/gobrailcreator"
)

// For a breadboard schematic refer to docs/images/PCA9501_Lamps_Buttons.png
// Just substidude the magnets with LED's and a 150Ohm resistor.

func main() {

	adatype, _ := gobrailcreator.ParseAdaptorType("digispark")
	gobrailcreator.Create(false, "dummy rob name",adatype , "./test/data/plan1.json", "./test/data/device_button4.json", "./test/data/device_togglebutton5.json")
}
