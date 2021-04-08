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
	loopCounter := 0
	var button *raildevices.ButtonDevice
	var togButton *raildevices.ToggleButtonDevice
	var lamp1 *raildevices.LampDevice
	var lamp2 *raildevices.LampDevice

	work := func() {
		gobot.Every(300*time.Millisecond, func() {
			if loopCounter == 0 {
				time.Sleep(2000 * time.Millisecond)
				fmt.Printf("\n------ Init Button ------\n")
				button, _ = raildevices.NewButton(boardAPI, boardID, 4, "Taste 1")
				togButton, _ = raildevices.NewToggleButton(boardAPI, boardID, 5, "Taste 2")
				lamp1, _ = raildevices.NewLamp(boardAPI, boardID, 0, "Strassenlampe 1", raildevices.Timing{})
				lamp2, _ = raildevices.NewLamp(boardAPI, boardID, 1, "Strassenlampe 2", raildevices.Timing{})
				fmt.Printf("\n------ Mapped pins ------\n")
				mPins := boardAPI.GetMappedAPIBinaryPins(boardID)
				fmt.Println(mPins)
				mPins = boardAPI.GetMappedAPIMemoryPins(boardID)
				fmt.Println(mPins)
				fmt.Printf("\n------ Map inputs ------\n")
				lamp1.Map(button)
				lamp2.Map(togButton)
				fmt.Printf("\n------ Now running ------\n")
			}
			//
			lamp1.Run()
			lamp2.Run()
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
