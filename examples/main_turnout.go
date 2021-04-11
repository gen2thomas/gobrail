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

var boardAPI *boardsapi.BoardsAPI

func main() {
	adaptor := digispark.NewAdaptor()
	boardAPI = boardsapi.NewBoardsAPI(adaptor)
	boardAPI.AddBoard(boardRecipePca9501)
	// setup IO's
	fmt.Printf("\n------ Init Inputs ------\n")
	togButton := createToggleButton("Taste 1", boardID, 4)
	fmt.Printf("\n------ Init Outputs ------\n")
	turnout := createTurnout("Weiche 1", boardID, 0, 3, raildevices.Timing{Starting: 1 * time.Second, Stopping: 1 * time.Second})
	lampRed := createLamp("Signal rot", boardID, 1, raildevices.Timing{Stopping: 50 * time.Millisecond})
	lampGreen := createLamp("Signal grün", boardID, 2, raildevices.Timing{Starting: 500 * time.Millisecond})
	fmt.Printf("\n------ Connect inputs to outputs ------\n")
	turnout.Connect(togButton)
	lampRed.Connect(turnout)
	lampGreen.ConnectInverse(lampRed)
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

func createToggleButton(railDeviceName string, boardID string, boardPinNr uint8) (toggleButton *raildevices.ToggleButtonDevice) {
	input, _ := boardAPI.GetInputPin(boardID, boardPinNr)
	toggleButton = raildevices.NewToggleButton(input, railDeviceName)
	return
}

func createLamp(railDeviceName string, boardID string, boardPinNr uint8, timing raildevices.Timing) (lamp *raildevices.LampDevice) {
	co := raildevices.NewCommonOutput(railDeviceName, timing, "lamp")
	output, _ := boardAPI.GetOutputPin(boardID, boardPinNr)
	lamp = raildevices.NewLamp(co, output)
	return
}

func createTurnout(railDeviceName string, boardID string, boardPinNrBranch uint8, boardPinNrMain uint8, timing raildevices.Timing) (turnout *raildevices.TurnoutDevice) {
	outputBranch, _ := boardAPI.GetOutputPin(boardID, boardPinNrBranch)
	outputMain, _ := boardAPI.GetOutputPin(boardID, boardPinNrMain)
	co := raildevices.NewCommonOutput(railDeviceName, timing, "turnout")
	turnout = raildevices.NewTurnout(co, outputBranch, outputMain)
	return
}