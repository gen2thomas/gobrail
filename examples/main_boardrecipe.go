// +build example
//
// Do not build by default.

package main

import (
	"fmt"

	"github.com/gen2thomas/gobrail/internal/boardrecipe"
)

func main() {
	var err error
	var recipe boardrecipe.Ingredients
	
  recipe, err = boardrecipe.ReadIngredients("./test/data/board_typ2_0x04.json")
  if err != nil {
  	fmt.Println("an error:", err)
  	return
  }
  fmt.Printf("Recipe - %s\n", recipe)
	
	if err = recipe.Verify(); err != nil{
		fmt.Printf("An error at '%s': %s\n", recipe.Name, err)
	}
  
  recipe, err = boardrecipe.ReadIngredients("./test/data/board_typ2_0x05.json")
  if err != nil {
  	fmt.Println("an error:", err)
  	return
  }
  fmt.Printf("Recipe - %s\n", recipe)
  if err = recipe.Verify(); err != nil{
		fmt.Printf("An error at '%s': %s\n", recipe.Name, err)
	}
}