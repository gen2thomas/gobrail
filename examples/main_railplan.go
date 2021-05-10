// +build example
//
// Do not build by default.

package main

import (
	"fmt"

	"github.com/gen2thomas/gobrail/internal/railplan"
)

const cookbook = "./test/data/plans/plan2.json"

func main() {
	var err error
	var book railplan.CookBook

	book, err = railplan.ReadCookBook(cookbook)
	if err != nil {
		fmt.Println("an error:", err)
		return
	}
	fmt.Printf("Cook book %s contains %d Board recipes and %d device recipes\n", cookbook, len(book.BoardRecipes), len(book.DeviceRecipes))

	for _, recipe := range book.BoardRecipes {
		fmt.Printf("Board ingredients - %s\n", recipe)
	}
	for _, recipe := range book.DeviceRecipes {
		fmt.Printf("Device ingredients - %s\n", recipe)
	}
}
