// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/digispark"

	"github.com/gen2thomas/gobrail/internal/boardsapi"
	"github.com/gen2thomas/gobrail/internal/raildevices"
)

// A toggle button is used to change the state of the railroad switch.
// Two separate lamps are used for simulate a red/green signal
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
	togButton, _ := raildevices.NewToggleButton(boardAPI, boardID, 4, "Taste 1")
	fmt.Printf("\n------ Init Outputs ------\n")
	turnout, _ := raildevices.NewTurnout(boardAPI, boardID, 0, "Weiche 1", 3, raildevices.Timing{Starting: 1 * time.Second, Stopping: 1 * time.Second})
	lampRed, _ := raildevices.NewLamp(boardAPI, boardID, 1, "Signal rot", raildevices.Timing{Stopping: 50 * time.Millisecond})
	lampGreen, _ := raildevices.NewLamp(boardAPI, boardID, 2, "Signal grün", raildevices.Timing{Starting: 500 * time.Millisecond})
	fmt.Printf("\n------ Map inputs to outputs ------\n")
	turnout.Connect(togButton)
	lampGreen.Connect(turnout)
	lampRed.ConnectInverse(lampGreen)
	fmt.Printf("\n------ Now running ------\n")

	work := func() {
		gobot.Every(50*time.Millisecond, func() {
			lampRed.Run()
			lampGreen.Run()
			turnout.Run()
		})
	}

	robot := gobot.NewRobot("play with button, turnout and lamps",
		[]gobot.Connection{adaptor},
		boardAPI.GobotDevices(),
		work,
	)

	err := robot.Start()
	if err != nil {
		fmt.Println(err)
	}
}
