// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/digispark"

	"github.com/gen2thomas/gobrail/internal/app/gobrailcreator"
)

// For a breadboard schematic refer to docs/images/PCA9501_Lamps_Buttons.png
// Just substidude the magnets with LED's and a 150Ohm resistor.

func main() {

	adaptor := digispark.NewAdaptor()

	// setup IO's
	gobrailcreator.Create(adaptor)	
	fmt.Printf("\n------ Now running ------\n")

	work := func() {
		gobot.Every(50*time.Millisecond, func() {
			gobrailcreator.Run()
		})
	}

	robot := gobot.NewRobot("play with button and lamp",
		[]gobot.Connection{adaptor},
		gobrailcreator.GobotDevices(),
		work,
	)

	err := robot.Start()
	if err != nil {
		fmt.Println(err)
	}
}
