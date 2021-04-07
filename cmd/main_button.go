package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/digispark"

	"github.com/gen2thomas/gobrail/internal/boardsapi"
	"github.com/gen2thomas/gobrail/internal/raildevices"
)

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
	var oldButtonState bool
	var button *raildevices.ButtonDevice
	var lamp *raildevices.LampDevice

	work := func() {
		gobot.Every(300*time.Millisecond, func() {
			if loopCounter == 0 {
				time.Sleep(2000 * time.Millisecond)
				fmt.Printf("\n------ Init Button ------\n")
				button, _ = raildevices.NewButton(boardAPI, boardID, 4, "Taste 1")
				lamp, _ = raildevices.NewLamp(boardAPI, boardID, 0, "Strassenlampe 1", raildevices.Timing{})
				fmt.Printf("\n------ Now running ------\n")
				fmt.Printf("\n------ Mapped pins ------\n")
				mPins := boardAPI.GetMappedAPIBinaryPins(boardID)
				fmt.Println(mPins)
				mPins = boardAPI.GetMappedAPIMemoryPins(boardID)
				fmt.Println(mPins)
			}
			buttonPressed, _ := button.IsPressed()
			if buttonPressed != oldButtonState {
				if buttonPressed {
					fmt.Printf("Button '%s' was pressed\n", button.Name())
					lamp.SwitchOn()
				} else {
					fmt.Printf("Button '%s' released\n", button.Name())
					lamp.SwitchOff()
				}
			}
			oldButtonState = buttonPressed
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
