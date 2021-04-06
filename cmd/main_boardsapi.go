package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/digispark"

	"github.com/gen2thomas/gobrail/internal/board"
	"github.com/gen2thomas/gobrail/internal/boardsapi"
)

const boardName = "IO_Mem_PCA9501"

var boardRecipePca9501 = boardsapi.BoardRecipe{
	Name:        boardName,
	ChipDevAddr: 0x04,
	BoardType:   boardsapi.Typ2,
}

var deviceArray = [...]string{
	"Weiche1 Links",
	"Weiche1 Rechts",
	"Weiche2 Links",
	"Weiche2 Rechts",
	"Signal1 Rot",
	"Signal1 Grün",
	"Signal2 Rot",
	"Signal2 Grün",
}

func main() {

	adaptor := digispark.NewAdaptor()
	boardAPI := boardsapi.NewBoardsAPI(adaptor, []boardsapi.BoardRecipe{boardRecipePca9501})
	firstLoop := true
	deviceArrayIdx := 0
	value := uint8(0)

	work := func() {
		gobot.Every(1000*time.Millisecond, func() {
			if firstLoop {
				fmt.Printf("\n------ IO test ------\n")
				boardAPI.SetAllOutputValues()
				time.Sleep(2000 * time.Millisecond)
				boardAPI.ResetAllOutputValues()
				time.Sleep(2000 * time.Millisecond)

				fmt.Printf("\n------ Free pins ------\n")
				freeAPIPins := boardAPI.GetFreeAPIPins(boardName, board.Binary)
				fmt.Println(freeAPIPins)

				fmt.Printf("\n------ Map pins ------\n")
				boardAPI.MapPin(boardName, 0, "Weiche1 Links")
				boardAPI.MapPin(boardName, 1, "Weiche1 Rechts")
				boardAPI.MapPin(boardName, 2, "Weiche2 Links")
				boardAPI.MapPin(boardName, 3, "Weiche2 Rechts")
				boardAPI.MapPin(boardName, 4, "Signal1 Rot")
				boardAPI.MapPin(boardName, 5, "Signal1 Grün")
				boardAPI.MapPin(boardName, 6, "Signal2 Rot")
				boardAPI.MapPin(boardName, 7, "Signal2 Grün")
				mappedAPIPins := boardAPI.GetMappedAPIPins(boardName, board.Binary)
				fmt.Println(mappedAPIPins)
				time.Sleep(2000 * time.Millisecond)

				firstLoop = false
				fmt.Printf("\n------ Now running ------\n")
			} else {
				boardAPI.SetValue(deviceArray[deviceArrayIdx], value)
				deviceArrayIdx++
				if deviceArrayIdx > 7 {
					deviceArrayIdx = 0
					value++
				}
				if value > 1 {
					value = 0
				}
			}
		})
	}

	robot := gobot.NewRobot("rotatePinsI2c",
		[]gobot.Connection{adaptor},
		boardAPI.GobotDevices(),
		work,
	)

	err := robot.Start()
	if err != nil {
		fmt.Println(err)
	}
}
