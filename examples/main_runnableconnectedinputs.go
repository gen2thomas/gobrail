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
	"github.com/gen2thomas/gobrail/internal/raildevicesapi"
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
	BoardType:   boardrecipe.Typ2,
}

var boardAPI *boardsapi.BoardsAPI

func main() {
	adaptor := digispark.NewAdaptor()
	boardAPI = boardsapi.NewBoardsAPI(adaptor)
	boardAPI.AddBoard(boardRecipePca9501)
	// setup IO's
	fmt.Printf("\n------ Init Inputs ------\n")
	button := createButton("Taste 1", boardID, 4)
	togButton := createToggleButton( "Taste 2", boardID, 5)
	fmt.Printf("\n------ Init Outputs ------\n")
	lamp1 := createLamp("Strassenlampe 1", boardID, 0, raildevices.Timing{})
	lamp2 := createLamp("Strassenlampe 2", boardID, 1, raildevices.Timing{Starting: 500*time.Millisecond})
	lamp3 := createLamp("Strassenlampe 3", boardID, 2, raildevices.Timing{Starting: time.Second, Stopping: time.Second})
	fmt.Printf("\n------ Connect inputs to outputs ------\n")
	lamp1.Connect(button)
	lamp2.Connect(togButton)
	lamp3.ConnectInverse(lamp2) // lamp3 will be switched off after lamp2 is really on
	fmt.Printf("\n------ Now running ------\n")

	work := func() {
		gobot.Every(50*time.Millisecond, func() {
			lamp3.Run()
			lamp1.Run()
			lamp2.Run()
		})
	}

	robot := gobot.NewRobot("play with connected buttons and lamps",
		[]gobot.Connection{adaptor},
		boardAPI.GobotDevices(),
		work,
	)

	err := robot.Start()
	if err != nil {
		fmt.Println(err)
	}
}

func createButton(railDeviceName string, boardID string, boardPinNr uint8) (button *raildevices.ButtonDevice) {
	input, _ := boardAPI.GetInputPin(boardID, boardPinNr)
	button = raildevices.NewButton(input, railDeviceName)
	return
}

func createToggleButton(railDeviceName string, boardID string, boardPinNr uint8) (toggleButton *raildevices.ToggleButtonDevice) {
	input, _ := boardAPI.GetInputPin(boardID, boardPinNr)
	toggleButton = raildevices.NewToggleButton(input, railDeviceName)
	return
}

func createLamp(railDeviceName string, boardID string, boardPinNr uint8, timing raildevices.Timing) (rd *runableDevice) {
	co := raildevices.NewCommonOutput(railDeviceName, timing)
	output, _ := boardAPI.GetOutputPin(boardID, boardPinNr)
	lamp := raildevices.NewLamp(co, output)
	rd = newRunableDevice(lamp)
	return
}

type runableDevice struct {
	raildevicesapi.Runner
	connectedInput raildevicesapi.Inputer
	inputInversion bool
	firstRun       bool
}

func newRunableDevice(outDev raildevicesapi.Runner) *runableDevice {
	return &runableDevice{
		Runner:   outDev,
		firstRun: true,
	}
}

// Connect is connecting an input for use in Run()
func (o *runableDevice) Connect(inputDevice raildevicesapi.Inputer) (err error) {
	if o.connectedInput != nil {
		return fmt.Errorf("The '%s' is already connected to an input '%s'", o.RailDeviceName(), o.connectedInput.RailDeviceName())
	}
	if o.RailDeviceName() == inputDevice.RailDeviceName() {
		return fmt.Errorf("Circular mapping blocked for '%s'", o.RailDeviceName())
	}
	o.connectedInput = inputDevice
	return nil
}

// ConnectInverse is connecting an input for use in Run(), but with inversed action
func (o *runableDevice) ConnectInverse(inputDevice raildevicesapi.Inputer) (err error) {
	o.Connect(inputDevice)
	o.inputInversion = true
	return nil
}

// RunCommon is called in a loop and will make action dependant on the input device
func (o *runableDevice) Run() (err error) {
	if o.connectedInput == nil {
		return fmt.Errorf("The '%s' can't run, please map to an input first", o.RailDeviceName())
	}
	var changed bool
	if changed, err = o.connectedInput.StateChanged(o.RailDeviceName()); err != nil {
		return err
	}
	if !(changed || o.firstRun) {
		return
	}
	o.firstRun = false
	if o.connectedInput.IsOn() != o.inputInversion {
		err = o.SwitchOn()
	} else {
		err = o.SwitchOff()
	}
	return
}

// ReleaseInput is used to unmap
func (o *runableDevice) ReleaseInput() {
	o.connectedInput = nil
}

