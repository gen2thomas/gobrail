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

var boardAPI *boardsapi.BoardsAPI

func main() {
	adaptor := digispark.NewAdaptor()
	boardAPI = boardsapi.NewBoardsAPI(adaptor, []boardsapi.BoardRecipe{boardRecipePca9501})
	// setup IO's
	fmt.Printf("\n------ Init Inputs ------\n")
	button := createButton("Taste 1", boardID, 4)
	togButton := createToggleButton( "Taste 2", boardID, 5)
	fmt.Printf("\n------ Init Outputs ------\n")
	lamp1 := createLamp("Strassenlampe 1", boardID, 0, raildevices.Timing{})
	lamp2 := createLamp("Strassenlampe 2", boardID, 1, raildevices.Timing{Starting: 500*time.Millisecond})
	lamp3 := createLamp("Strassenlampe 3", boardID, 2, raildevices.Timing{Starting: time.Second, Stopping: time.Second})
	fmt.Printf("\n------ Connect inputs to outputs ------\n")
	lamp1.Connect(button)
	lamp2.Connect(togButton)
	lamp3.ConnectInverse(lamp2) // lamp3 will be switched off after lamp2 is really on
	fmt.Printf("\n------ Now running ------\n")

	work := func() {
		gobot.Every(50*time.Millisecond, func() {
			lamp3.Run()
			lamp1.Run()
			lamp2.Run()
		})
	}

	robot := gobot.NewRobot("play with connected buttons and lamps",
		[]gobot.Connection{adaptor},
		boardAPI.GobotDevices(),
		work,
	)

	err := robot.Start()
	if err != nil {
		fmt.Println(err)
	}
}

func createButton(railDeviceName string, boardID string, boardPinNr uint8) (button *raildevices.ButtonDevice) {
	input, _ := boardAPI.GetInputPin(boardID, boardPinNr)
	button = raildevices.NewButton(input, railDeviceName)
	return
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
