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

  recipe, err = boardrecipe.ReadIngredients("./test/data/boardrecipes/board_typ2_0x04.json")
  if err != nil {
    fmt.Println("an error:", err)
    return
  }
  fmt.Printf("Recipe - %s\n", recipe)

  recipe, err = boardrecipe.ReadIngredients("./test/data/boardrecipes/board_typ2_0x05.json")
  if err != nil {
    fmt.Println("an error:", err)
    return
  }
  fmt.Printf("Recipe - %s\n", recipe)
}
