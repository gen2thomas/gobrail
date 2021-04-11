package gobrailcreator

// * gets configuration from config/reader (creation plan)
// * can add elements (by json or objects?)
// * can provide creation to run railroad

import (
	"fmt"
	"strings"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"

	"github.com/gen2thomas/gobrail/internal/boardpin"
	"github.com/gen2thomas/gobrail/internal/boardsapi"
	"github.com/gen2thomas/gobrail/internal/raildevices"
)

// RailDevice can run in loops
type RailDevice interface {
	raildevices.Outputer
	raildevices.Inputer
}

// BoardsAPIer is an interface for interact with a boards API2
type BoardsAPIer interface {
	GobotDevices() []gobot.Device
	GetInputPin(boardID string, boardPinNr uint8) (boardPin *boardpin.Input, err error)
	GetOutputPin(boardID string, boardPinNr uint8) (boardPin *boardpin.Output, err error)
}

const boardID = "IO_Mem_PCA9501"

var boardRecipePca9501 = boardsapi.BoardRecipe{
	Name:        boardID,
	ChipDevAddr: 0x04,
	BoardType:   boardsapi.Typ2,
}

var usedDevices = make(map[string]struct{})
var runningDevices []raildevices.Outputer
var boardAPI BoardsAPIer

// Create will create a static device connection for run
func Create(adaptor i2c.Connector) (err error) {
	boardAPI = boardsapi.NewBoardsAPI(adaptor, []boardsapi.BoardRecipe{boardRecipePca9501})
	if boardAPI == nil {
		return fmt.Errorf("BoardsAPI can't created")
	}

	fmt.Printf("\n------ Init Inputs ------\n")
	togButton := createToggleButton("Taste 2", boardRecipePca9501.Name, 5)
	button := createButton("Taste 1", boardRecipePca9501.Name, 4)
	fmt.Printf("\n------ Init Outputs ------\n")
	lampY1 := createLamp("Signal rot", boardRecipePca9501.Name, 0, raildevices.Timing{Stopping: 50 * time.Millisecond})
	lampY2 := createLamp("Signal grün", boardRecipePca9501.Name, 3, raildevices.Timing{})
	signal := createTwoLightSignal("Rot grün Signal", boardRecipePca9501.Name, 2, 1, raildevices.Timing{})
	fmt.Printf("\n------ Map inputs to outputs ------\n")
	lampY1.Connect(button)
	lampY2.Connect(lampY1)
	signal.Connect(togButton)
	runningDevices = append(runningDevices, lampY1)
	runningDevices = append(runningDevices, lampY2)
	runningDevices = append(runningDevices, signal)
	return
}

// Run calls the run functions of all devices
func Run() {
	for _, runningDevice := range runningDevices {
		runningDevice.Run()
	}
}

// GobotDevices gets all gobot devices of all boards
func GobotDevices() []gobot.Device {
	return boardAPI.GobotDevices()
}

func getFreeKey(railDeviceName string) (railDeviceKey string, err error) {
	railDeviceKey = strings.Replace(strings.ToLower(railDeviceName), " ", "_", -1)
	if _, ok := usedDevices[railDeviceKey]; ok {
		return "", fmt.Errorf("Rail device '%s' (key: %s) already in use", railDeviceName, railDeviceKey)
	}
	return
}

func createButton(railDeviceName string, boardID string, boardPinNr uint8) (button raildevices.Inputer) {
	railDeviceKey, _ := getFreeKey(railDeviceName)
	input, _ := boardAPI.GetInputPin(boardID, boardPinNr)
	button = raildevices.NewButton(input, railDeviceName)
	usedDevices[railDeviceKey] = struct{}{}
	return
}

func createToggleButton(railDeviceName string, boardID string, boardPinNr uint8) (toggleButton raildevices.Inputer) {
	railDeviceKey, _ := getFreeKey(railDeviceName)
	input, _ := boardAPI.GetInputPin(boardID, boardPinNr)
	toggleButton = raildevices.NewToggleButton(input, railDeviceName)
	usedDevices[railDeviceKey] = struct{}{}
	return
}

func createLamp(railDeviceName string, boardID string, boardPinNr uint8, timing raildevices.Timing) (lamp RailDevice) {
	railDeviceKey, _ := getFreeKey(railDeviceName)
	co := raildevices.NewCommonOutput(railDeviceName, timing, "lamp")
	output, _ := boardAPI.GetOutputPin(boardID, boardPinNr)
	lamp = raildevices.NewLamp(co, output)
	usedDevices[railDeviceKey] = struct{}{}
	return
}

func createTwoLightSignal(railDeviceName string, boardID string, boardPinNrPass uint8, boardPinNrStop uint8, timing raildevices.Timing) (signal RailDevice) {
	railDeviceKey, _ := getFreeKey(railDeviceName)
	outputPass, _ := boardAPI.GetOutputPin(boardID, boardPinNrPass)
	outputStop, _ := boardAPI.GetOutputPin(boardID, boardPinNrStop)
	co := raildevices.NewCommonOutput(railDeviceName, timing, "two light signal")
	signal = raildevices.NewTwoLightsSignal(co, outputPass, outputStop)
	usedDevices[railDeviceKey] = struct{}{}
	return
}

func createTurnout(railDeviceName string, boardID string, boardPinNrBranch uint8, boardPinNrMain uint8, timing raildevices.Timing) (turnout RailDevice) {
	railDeviceKey, _ := getFreeKey(railDeviceName)
	outputBranch, _ := boardAPI.GetOutputPin(boardID, boardPinNrBranch)
	outputMain, _ := boardAPI.GetOutputPin(boardID, boardPinNrMain)
	co := raildevices.NewCommonOutput(railDeviceName, timing, "turnout")
	turnout = raildevices.NewTurnout(co, outputBranch, outputMain)
	usedDevices[railDeviceKey] = struct{}{}
	return
}
