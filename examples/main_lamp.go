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
	loopCounter := 0
	fmt.Printf("\n------ Init Lamp ------\n")
	lamp := createLamp("Strassenlampe 1", boardID, 0, raildevices.Timing{})
	fmt.Printf("\n------ Used pins ------\n")
	uPins := boardAPI.GetUsedPins(boardID)
	fmt.Println(uPins)
	fmt.Printf("\n------ Now running ------\n")

	work := func() {
		gobot.Every(4000*time.Millisecond, func() {
			lamp.SwitchOn()
			time.Sleep(2000 * time.Millisecond)
			lamp.SwitchOff()
			if loopCounter == 2 {
				if deferr := lamp.MakeDefective(); deferr == nil {
					fmt.Printf("Lamp '%s' is now defective, please repair\n", lamp.RailDeviceName())
				}
			}
			if loopCounter == 5 {
				if reperr := lamp.Repair(); reperr == nil {
					fmt.Printf("Lamp '%s' was repaired\n", lamp.RailDeviceName())
				}
			}
			if isDefectErr := lamp.IsDefective(); isDefectErr != nil {
				fmt.Printf("Lamp '%s' is defective\n", lamp.RailDeviceName())
			} else {
			  fmt.Printf("Lamp '%s' is working\n", lamp.RailDeviceName())
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

func createLamp(railDeviceName string, boardID string, boardPinNr uint8, timing raildevices.Timing) (lamp *raildevices.LampDevice) {
	co := raildevices.NewCommonOutput(railDeviceName, timing, "lamp")
	output, _ := boardAPI.GetOutputPin(boardID, boardPinNr)
	lamp = raildevices.NewLamp(co, output)
	return
}
