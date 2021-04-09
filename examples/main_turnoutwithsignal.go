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
// Afterwards a red/green signal is switched without any delay.
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
	//button, _ := raildevices.NewButton(boardAPI, boardID, 4, "Taste 1")
	togButton, _ := raildevices.NewToggleButton(boardAPI, boardID, 5, "Taste 2")
	fmt.Printf("\n------ Init Outputs ------\n")
	turnout, _ := raildevices.NewTurnout(boardAPI, boardID, 0, "Weiche 1", 3, raildevices.Timing{Starting: 500 * time.Millisecond, Stopping: 500 * time.Millisecond})
	redgreensignal, _ := raildevices.NewTwoLightSignal(boardAPI, boardID, 2, "Signal rot grün", 1, raildevices.Timing{})
	fmt.Printf("\n------ Map inputs to outputs ------\n")
	turnout.Connect(togButton)
	redgreensignal.Connect(turnout)
	fmt.Printf("\n------ Now running ------\n")

	work := func() {
		gobot.Every(50*time.Millisecond, func() {
			if err := redgreensignal.Run(); err != nil{
				fmt.Println(err)
			}
			if err := turnout.Run(); err != nil{
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("play with button, turnout and light signal",
		[]gobot.Connection{adaptor},
		boardAPI.GobotDevices(),
		work,
	)

	err := robot.Start()
	if err != nil {
		fmt.Println(err)
	}
}
