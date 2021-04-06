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
	var lamp *raildevices.LampDevice

	work := func() {
		gobot.Every(4000*time.Millisecond, func() {
			if loopCounter == 0 {
				fmt.Printf("\n------ Init Lamp ------\n")
				lamp = raildevices.NewLamp(boardAPI, boardID, 0, "Strassenlampe 1", raildevices.Timing{})
				fmt.Printf("\n------ Now running ------\n")
				fmt.Printf("\n------ Mapped pins ------\n")
				mPins := boardAPI.GetMappedAPIBinaryPins(boardID)
				fmt.Println(mPins)
				mPins = boardAPI.GetMappedAPIMemoryPins(boardID)
				fmt.Println(mPins)
			}
			lamp.SwitchOn()
			time.Sleep(2000 * time.Millisecond)
			lamp.SwitchOff()
			if loopCounter == 2 {
				lamp.MakeDefective()
			}
			if loopCounter == 5 {
				lamp.Repair()
			}
			loopCounter++
		})
	}

	robot := gobot.NewRobot("play with lamp",
		[]gobot.Connection{adaptor},
		boardAPI.GobotDevices(),
		work,
	)

	err := robot.Start()
	if err != nil {
		fmt.Println(err)
	}
}
