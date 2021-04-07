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
	var button *raildevices.ButtonDevice
	var lamp1 *raildevices.LampDevice
	var lamp2 *raildevices.LampDevice

	work := func() {
		gobot.Every(300*time.Millisecond, func() {
			if loopCounter == 0 {
				time.Sleep(2000 * time.Millisecond)
				fmt.Printf("\n------ Init Button ------\n")
				button, _ = raildevices.NewButton(boardAPI, boardID, 4, "Taste 1")
				lamp1, _ = raildevices.NewLamp(boardAPI, boardID, 0, "Strassenlampe 1", raildevices.Timing{})
				lamp2, _ = raildevices.NewLamp(boardAPI, boardID, 1, "Strassenlampe 2", raildevices.Timing{})
				fmt.Printf("\n------ Now running ------\n")
				fmt.Printf("\n------ Mapped pins ------\n")
				mPins := boardAPI.GetMappedAPIBinaryPins(boardID)
				fmt.Println(mPins)
				mPins = boardAPI.GetMappedAPIMemoryPins(boardID)
				fmt.Println(mPins)
				lamp1.SwitchOff()
				lamp2.SwitchOff()
			}
			buttonPressed, _ := button.IsPressed()
			if buttonPressed {
				lamp1.SwitchOn()
			} else {
				lamp1.SwitchOff()
			}
			buttonChanged, _ := button.IsChanged()
			if buttonChanged {
				if button.WasPressed() {
					fmt.Printf("Button '%s' was pressed\n", button.Name())
					lamp2.SwitchOn()
				} else {
					fmt.Printf("Button '%s' released\n", button.Name())
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
