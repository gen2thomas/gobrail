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

func main() {

	adaptor := digispark.NewAdaptor()
	boardAPI := boardsapi.NewBoardsAPI(adaptor, []boardsapi.BoardRecipe{boardRecipePca9501})
	firstLoop := true

	work := func() {
		gobot.Every(4000*time.Millisecond, func() {
			if firstLoop {
				//boardAPI.ShowConfigs()
				fmt.Printf("\n------ Free pins ------\n")
				freeAPIPins := boardAPI.GetFreeAPIPins(boardName, board.Binary)
				fmt.Println(freeAPIPins)

				fmt.Printf("\n------ Mapped pins ------\n")
				mappedAPIPins := boardAPI.GetMappedAPIPins(boardName, board.Binary)
				fmt.Println(mappedAPIPins)

				fmt.Printf("\n------ Map pin ------\n")
				boardAPI.MapPin(boardName, 0, "Weiche1 Links")
				mappedAPIPins = boardAPI.GetMappedAPIPins(boardName, board.Binary)
				fmt.Println(mappedAPIPins)

				// already mapped
				boardAPI.MapPin(boardName, 0, "Weiche1 Links")
				boardAPI.MapPin(boardName, 0, "Weiche1 Rechts")

				boardAPI.SetValue("Weiche1 Links", 0)

				// not mapped
				boardAPI.SetValue("Weiche1 Rechts", 0)

				fmt.Printf("\n------ Release pin ------\n")
				boardAPI.ReleasePin("Weiche1 Links")
				boardAPI.SetValue("Weiche1 Links", 0)

				// already released
				boardAPI.ReleasePin("Weiche1 Links")
				//fmt.Printf("\n------ Write to Memory ------\n")
				//boardAPI.WriteBoardConfig()
				//fmt.Printf("\n------ Read from Memory ------\n")
				//boardAPI.ReadBoardConfig()
				//boardAPI.ShowConfigs()
				firstLoop = false
				fmt.Printf("\n------ Now running ------\n")
			}

			//boardAPI.SetAllOutputValues()
			time.Sleep(2000 * time.Millisecond)
			//boardAPI.ResetAllOutputValues()
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
