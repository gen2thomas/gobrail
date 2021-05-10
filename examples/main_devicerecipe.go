// +build example
//
// Do not build by default.

package main

import (
	"fmt"

	"github.com/gen2thomas/gobrail/internal/devicerecipe"
)

func main() {
	var err error
	var recipe devicerecipe.Ingredients

	recipe, err = devicerecipe.ReadIngredients("./test/data/device_button4.json")
	if err != nil {
		fmt.Println("an error:", err)
		return
	}
	fmt.Printf("Device - %s\n", recipe)

	recipe, err = devicerecipe.ReadIngredients("./test/data/device_togglebutton5.json")
	if err != nil {
		fmt.Println("an error:", err)
		return
	}
	fmt.Printf("Device - %s\n", recipe)
}
