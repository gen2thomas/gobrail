package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/digispark"

	"github.com/gen2thomas/gobrail/internal/boardsapi"
	"github.com/gen2thomas/gobrail/internal/raildevices"
)

// Two buttons are used to switch on and off a lamp.
// First button is used as normal button.
// Second button is used as toggle button.
//
// For a breadboard schematic refer to docs/images/PCA9501_Lamps_Buttons.png
// Just substidude the magnets with LED's and a 150Ohm resistor.

const boardID = "IO_Mem_PCA9501"

var boardRecipePca9501 = boardsapi.BoardRecipe{
	Name:        boardID,
	ChipDevAddr: 0x04,
	BoardType:   boardsapi.Typ2,
}

func main() {

	adaptor := digispark.NewAdaptor()
	boardAPI := boardsapi.NewBoardsAPI(adaptor, []boardsapi.BoardRecipe{boardRecipePca9501})
	// setup IO's
	fmt.Printf("\n------ Init Inputs ------\n")
	button, _ := raildevices.NewButton(boardAPI, boardID, 4, "Taste 1")
	togButton, _ := raildevices.NewToggleButton(boardAPI, boardID, 5, "Taste 2")
	fmt.Printf("\n------ Init Outputs ------\n")
	lamp1, _ := raildevices.NewLamp(boardAPI, boardID, 0, "Strassenlampe 1", raildevices.Timing{})
	lamp2, _ := raildevices.NewLamp(boardAPI, boardID, 1, "Strassenlampe 2", raildevices.Timing{})
	lamp3, _ := raildevices.NewLamp(boardAPI, boardID, 2, "Strassenlampe 3", raildevices.Timing{})
	fmt.Printf("\n------ Map inputs to outputs ------\n")
	lamp1.Map(button)
	lamp2.Map(togButton)
	lamp3.Map(lamp2) // lamp3 will be switched on after lamp2 is really on
	fmt.Printf("\n------ Now running ------\n")

	work := func() {
		gobot.Every(50*time.Millisecond, func() {
			lamp3.Run()
			lamp1.Run()
			lamp2.Run()
		})
	}

	robot := gobot.NewRobot("play with button and lamp",
		[]gobot.Connection{adaptor},
		boardAPI.GobotDevices(),
		work,
	)

	err := robot.Start()
	if err != nil {
		fmt.Println(err)
	}
}
