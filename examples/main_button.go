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
	"github.com/gen2thomas/gobrail/internal/boardrecipe"
)

// Two buttons are used to switch on and off a lamp.
// First button is used as normal button.
// Second button is used as toggle button.
//
// For a breadboard schematic refer to docs/images/PCA9501_Lamps_Buttons.png
// Just substitute the magnets with LED's and a 150Ohm resistor.

const boardID = "IO_Mem_PCA9501"

var boardRecipePca9501 = boardrecipe.Ingredients{
	Name:        boardID,
	ChipDevAddr: 0x04,
	Type:   "Type2",
}

var boardAPI *boardsapi.BoardsAPI

func main() {

	adaptor := digispark.NewAdaptor()
	boardAPI = boardsapi.NewBoardsAPI(adaptor)
	boardAPI.AddBoard(boardRecipePca9501)
	loopCounter := 0
	fmt.Printf("\n------ Init Button ------\n")
	button := createButton("Taste 1", boardID, 4)
	togButton := createToggleButton("Taste 2", boardID, 5)
	lamp1 := createLamp("Strassenlampe 1", boardID, 0 ,raildevices.Timing{})
	lamp2 := createLamp("Strassenlampe 2", boardID, 1 ,raildevices.Timing{})
	fmt.Printf("\n------ Used pins ------\n")
	uPins := boardAPI.GetUsedPins(boardID)
	fmt.Println(uPins)
	fmt.Printf("\n------ Now running ------\n")

	work := func() {
		gobot.Every(50*time.Millisecond, func() {
			//
			if changed, _ := button.StateChanged("v"); changed {
				if button.IsOn() {
					fmt.Printf("Button '%s' was pressed\n", button.RailDeviceName())
					lamp1.SwitchOn()
				} else {
					fmt.Printf("Button '%s' released\n", button.RailDeviceName())
					lamp1.SwitchOff()
				}
			}
			//
			if changed, _ := togButton.StateChanged("v"); changed {
				if togButton.IsOn() {
					fmt.Printf("Toggle '%s' to on\n", togButton.RailDeviceName())
					lamp2.SwitchOn()
				} else {
					fmt.Printf("Toggle '%s' to off\n", togButton.RailDeviceName())
					lamp2.SwitchOff()
				}
			}
			loopCounter++
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
	co := raildevices.NewCommonOutput(railDeviceName, timing)
	output, _ := boardAPI.GetOutputPin(boardID, boardPinNr)
	lamp = raildevices.NewLamp(co, output)
	return
}
