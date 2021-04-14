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
	var deviceRecipe devicerecipe.RailDeviceRecipe
	
  deviceRecipes, err := devicerecipe.ReadPlan("./test/data/plan1.json")
  if err != nil {
  	fmt.Println("an error:", err)
  }

  if err == nil {
		fmt.Printf("Now Print %d Recipes:\n", len(deviceRecipes))
		for _, deviceRecipe := range deviceRecipes {
			fmt.Println(deviceRecipe)
		}
	}
  
	deviceRecipe, err = devicerecipe.ReadDevice("./test/data/device_button4.json")
  if err != nil {
  	fmt.Println("an error:", err)
  }
  if err == nil {
		fmt.Printf("Now Print Recipe '%s':\n", deviceRecipe.Name)
		fmt.Println(deviceRecipe)
	}
	if err = deviceRecipe.Verify(); err != nil{
		fmt.Printf("An error at '%s': %s\n", deviceRecipe.Name, err)
	}
  
  deviceRecipe, err = devicerecipe.ReadDevice("./test/data/device_togglebutton5.json")
  if err != nil {
  	fmt.Println("an error:", err)
  }
  if err == nil {
		fmt.Printf("Now Print Recipe '%s':\n", deviceRecipe.Name)
		fmt.Println(deviceRecipe)
	}
}