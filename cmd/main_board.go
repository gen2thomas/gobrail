package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/digispark"

	"github.com/gen2thomas/gobrail/internal/board"
)

// To test this program, a simple board is used with a PCA9501 and some standard LED's (2V, 20mA).
// Catode of LED is connected to chip IO, anode with 150Ohm resistors to +5V ("low active").
// SDA & SCL are connected to platforms related output. For digispark each with 10kOhm pullup to +5V.
//
// When using "low active" mode the the total power disapation is lower than allowed 400mW.
// --> 0,4V*20mA = 8mW --> 8Pins ==> 64mW
// Please consider maximum ratings when using "high active" mode.
//
// It is possible to use another platform than digispark. Some has the pullup resistors already in place.

type boardRecipe struct {
	Name        string
	ChipDevAddr uint8
}

const boardName = "IO_Mem_PCA9501"

var boardRecipePca9501 = boardRecipe{
	Name:        boardName,
	ChipDevAddr: 0x04,
}

func main() {

	adaptor := digispark.NewAdaptor()
	board := board.NewBoardTyp2(adaptor, boardRecipePca9501.ChipDevAddr, boardRecipePca9501.Name)
	firstLoop := true
	pin := uint8(0)
	value := uint8(0)

	work := func() {
		gobot.Every(1000*time.Millisecond, func() {
			if firstLoop {
				fmt.Printf("\n------ Config ------\n")
				board.ShowBoardConfig()
				time.Sleep(500 * time.Millisecond)

				fmt.Printf("\n------ IO test ------\n")
				board.SetAllIoPins()
				time.Sleep(2000 * time.Millisecond)
				board.ResetAllIoPins()
				time.Sleep(2000 * time.Millisecond)

				fmt.Printf("\n------ Now running ------\n")
				firstLoop = false
			} else {
				board.SetValue(pin, value)
				pin++
				if pin > 7 {
					pin = 0
					value++
				}
				if value > 1 {
					value = 0
				}
			}
		})
	}

	robot := gobot.NewRobot("try board TYP2",
		[]gobot.Connection{adaptor},
		board.Devices(),
		work,
	)

	err := robot.Start()
	if err != nil {
		fmt.Println(err)
	}
}
