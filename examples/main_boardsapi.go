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
	"github.com/gen2thomas/gobrail/internal/boardpin"
	"github.com/gen2thomas/gobrail/internal/boardrecipe"
)

const boardName = "IO_Mem_PCA9501"

var boardRecipePca9501 = boardsapi.BoardRecipe{
	Name:        boardName,
	ChipDevAddr: 0x04,
	BoardType:   boardrecipe.Typ2,
}

var deviceArray [4]boardpin.Output

func main() {

	adaptor := digispark.NewAdaptor()
	boardAPI := boardsapi.NewBoardsAPI(adaptor)
	boardAPI.AddBoard(boardRecipePca9501)
	deviceArrayIdx := 0
	value := uint8(0)
	fmt.Printf("\n------ Free pins ------\n")
	freePins := boardAPI.GetFreePins(boardName)
	fmt.Println(freePins)
	fmt.Printf("\n------ Map pins ------\n")
	weiche1Links,_ := boardAPI.GetOutputPin(boardName, 0)
	weiche1Rechts,_:= boardAPI.GetOutputPin(boardName, 3)
	signal1Rot,_:=boardAPI.GetOutputPin(boardName, 1)
	signal1Gruen,_:= boardAPI.GetOutputPin(boardName, 2)
	usedPins := boardAPI.GetUsedPins(boardName)
	fmt.Println(usedPins)
  deviceArray[0] = *weiche1Links
  deviceArray[1] = *signal1Rot
  deviceArray[2] = *weiche1Rechts
  deviceArray[3] = *signal1Gruen
	fmt.Printf("\n------ Now running ------\n")

	work := func() {
		gobot.Every(500*time.Millisecond, func() {
			deviceArray[deviceArrayIdx].WriteValue(value)
			deviceArrayIdx++
			if deviceArrayIdx > 3 {
				deviceArrayIdx = 0
				value++
			}
			if value > 1 {
				value = 0
			}
		})
	}

	robot := gobot.NewRobot("rotate Pins",
		[]gobot.Connection{adaptor},
		boardAPI.GobotDevices(),
		work,
	)

	err := robot.Start()
	if err != nil {
		fmt.Println(err)
	}
}
