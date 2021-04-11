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

var boardAPI *boardsapi.BoardsAPI

func main() {

	adaptor := digispark.NewAdaptor()
	boardAPI = boardsapi.NewBoardsAPI(adaptor, []boardsapi.BoardRecipe{boardRecipePca9501})
	// setup IO's
	fmt.Printf("\n------ Init Inputs ------\n")
	togButton := createToggleButton("Taste 2", boardID, 5)
	fmt.Printf("\n------ Init Outputs ------\n")
	turnout := createTurnout("Weiche 1", boardID, 0, 3, raildevices.Timing{Starting: 500 * time.Millisecond, Stopping: 500 * time.Millisecond})
	redgreensignal := createTwoLightSignal("Signal rot grün", boardID, 1, 2, raildevices.Timing{})
	fmt.Printf("\n------ Connect inputs to outputs ------\n")
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

func createToggleButton(railDeviceName string, boardID string, boardPinNr uint8) (toggleButton *raildevices.ToggleButtonDevice) {
	input, _ := boardAPI.GetInputPin(boardID, boardPinNr)
	toggleButton = raildevices.NewToggleButton(input, railDeviceName)
	return
}

func createTurnout(railDeviceName string, boardID string, boardPinNrBranch uint8, boardPinNrMain uint8, timing raildevices.Timing) (turnout *raildevices.TurnoutDevice) {
	outputBranch, _ := boardAPI.GetOutputPin(boardID, boardPinNrBranch)
	outputMain, _ := boardAPI.GetOutputPin(boardID, boardPinNrMain)
	co := raildevices.NewCommonOutput(railDeviceName, timing, "turnout")
	turnout = raildevices.NewTurnout(co, outputBranch, outputMain)
	return
}

func createTwoLightSignal(railDeviceName string, boardID string, boardPinNrPass uint8, boardPinNrStop uint8, timing raildevices.Timing) (signal *raildevices.TwoLightsSignalDevice) {
	outputPass, _ := boardAPI.GetOutputPin(boardID, boardPinNrPass)
	outputStop, _ := boardAPI.GetOutputPin(boardID, boardPinNrStop)
	co := raildevices.NewCommonOutput(railDeviceName, timing, "two light signal")
	signal = raildevices.NewTwoLightsSignal(co, outputPass, outputStop)
	return
}
