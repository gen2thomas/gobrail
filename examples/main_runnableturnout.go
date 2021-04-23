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
	"github.com/gen2thomas/gobrail/internal/boardrecipe"
)

// A toggle button is used to change the state of the railroad switch.
// Two separate lamps are used for simulate a red/green signal
//
// For a breadboard schematic refer to docs/images/PCA9501_Lamps_Buttons.png
// Just substitute the magnets with LED's and a 150Ohm resistor.

const boardID = "IO_Mem_PCA9501"

var boardRecipePca9501 = boardrecipe.Ingredients{
	Name:        boardID,
	ChipDevAddr: 0x04,
	Type:   "Type2",
}

var boardAPI *boardsapi.BoardsAPI

func main() {
	adaptor := digispark.NewAdaptor()
	boardAPI = boardsapi.NewBoardsAPI(adaptor)
	boardAPI.AddBoard(boardRecipePca9501)
	// setup IO's
	fmt.Printf("\n------ Init Inputs ------\n")
	togButton := createToggleButton("Taste 1", boardID, 4)
	fmt.Printf("\n------ Init Outputs ------\n")
	turnout := createTurnout("Weiche 1", boardID, 0, 3, raildevices.Timing{Starting: 1 * time.Second, Stopping: 1 * time.Second})
	lampRed := createLamp("Signal rot", boardID, 1, raildevices.Timing{Stopping: 50 * time.Millisecond})
	lampGreen := createLamp("Signal grün", boardID, 2, raildevices.Timing{Starting: 500 * time.Millisecond})
	fmt.Printf("\n------ Connect inputs to outputs ------\n")
	turnout.Connect(togButton)
	lampRed.Connect(turnout)
	lampGreen.ConnectInverse(lampRed)
	fmt.Printf("\n------ Now running ------\n")

	work := func() {
		gobot.Every(50*time.Millisecond, func() {
			lampRed.Run()
			lampGreen.Run()
			turnout.Run()
		})
	}

	robot := gobot.NewRobot("play with button, turnout and lamps",
		[]gobot.Connection{adaptor},
		boardAPI.GobotDevices(),
		work,
	)

	err := robot.Start()
	if err != nil {
		fmt.Println(err)
	}
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

func createTurnout(railDeviceName string, boardID string, boardPinNrBranch uint8, boardPinNrMain uint8, timing raildevices.Timing) (rd *runableDevice) {
	outputBranch, _ := boardAPI.GetOutputPin(boardID, boardPinNrBranch)
	outputMain, _ := boardAPI.GetOutputPin(boardID, boardPinNrMain)
	co := raildevices.NewCommonOutput(railDeviceName, timing)
	turnout := raildevices.NewTurnout(co, outputBranch, outputMain)
	rd = newRunableDevice(turnout)
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

// RunCommon is called in a loop and will make action dependent on the input device
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