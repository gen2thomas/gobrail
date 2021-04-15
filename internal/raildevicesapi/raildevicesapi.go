package raildevicesapi

import (
	"fmt"
	"strings"
	"time"

	"github.com/gen2thomas/gobrail/internal/boardpin"
	"github.com/gen2thomas/gobrail/internal/devicerecipe"
	"github.com/gen2thomas/gobrail/internal/raildevices"
)

// Inputer is an interface for input devices to map in output devices. When an output device
// have this functions it can be used as input for an successive device.
type Inputer interface {
	RailDeviceName() string
	StateChanged(visitor string) (hasChanged bool, err error)
	IsOn() bool
}

// Outputer is an interface for output devices
type Outputer interface {
	RailDeviceName() string
	SwitchOn() (err error)
	SwitchOff() (err error)
}

// Runner is an interface for devices which can call cyclic
type Runner interface {
	Inputer
	Outputer
}

// BoardsIOAPIer is an interface for interact with a boards API which provides IO pins
type BoardsIOAPIer interface {
	GetInputPin(boardID string, boardPinNr uint8) (boardPin *boardpin.Input, err error)
	GetOutputPin(boardID string, boardPinNr uint8) (boardPin *boardpin.Output, err error)
}

// RailDeviceAPI describes the api
type RailDeviceAPI struct {
	boardsIOAPI    BoardsIOAPIer
	devices        map[string]struct{}
	runableDevices map[string]*runableDevice
	inputDevices   map[string]Inputer
	connections    map[string]string
}

// NewRailDevicesAPI creates a new instance of rail device API
func NewRailDevicesAPI(boardsIOAPI BoardsIOAPIer) *RailDeviceAPI {
	return &RailDeviceAPI{
		devices:        make(map[string]struct{}),
		boardsIOAPI:    boardsIOAPI,
		runableDevices: make(map[string]*runableDevice),
		inputDevices:   make(map[string]Inputer),
		connections:    make(map[string]string),
	}
}

// AddDevice creates a device from recipe and add it to the list
func (di *RailDeviceAPI) AddDevice(deviceRecipe devicerecipe.RailDeviceRecipe) (err error) {
	railDeviceKey := getKey(deviceRecipe.Name)
	if _, ok := di.devices[railDeviceKey]; ok {
		return fmt.Errorf("Rail device '%s' (key: %s) already in use", deviceRecipe.Name, railDeviceKey)
	}
	var inDev Inputer
	var runDev *runableDevice
	switch devicerecipe.TypeMap[deviceRecipe.Type] {
	case devicerecipe.Button:
		inDev = di.createButton(deviceRecipe)
	case devicerecipe.ToggleButton:
		inDev = di.createToggleButton(deviceRecipe)
	case devicerecipe.Lamp:
		runDev = di.createLamp(deviceRecipe)
	case devicerecipe.TwoLightsSignal:
		runDev = di.createTwoLightSignal(deviceRecipe)
	case devicerecipe.Turnout:
		runDev = di.createTurnout(deviceRecipe)
	default:
		return fmt.Errorf("Unknown type '%s'", deviceRecipe.Type)
	}
	if inDev != nil {
		di.inputDevices[railDeviceKey] = inDev
	}
	if runDev != nil {
		di.runableDevices[railDeviceKey] = runDev
	}
	if deviceRecipe.Connect != "" {
		di.connections[railDeviceKey] = getKey(deviceRecipe.Connect)
	}
	di.devices[railDeviceKey] = struct{}{}
	return
}

// ConnectNow create all connections
func (di *RailDeviceAPI) ConnectNow() (err error) {
	for runningDevKey, runableDevice := range di.runableDevices {
		var conKey string
		var ok bool
		if conKey, ok = di.connections[runningDevKey]; !ok {
			continue
		}
		var conDev Inputer
		if conDev, ok = di.runableDevices[conKey]; !ok {
			conDev = di.inputDevices[conKey]
		}
		if conDev != nil {
			runableDevice.Connect(conDev)
		} else {
			return fmt.Errorf("Device with key '%s' to connect with '%s' not found", conKey, runableDevice.RailDeviceName())
		}
	}
	return
}

// Run calls the run functions of all runnable devices
func (di *RailDeviceAPI) Run() {
	for _, runableDevice := range di.runableDevices {
		runableDevice.Run()
	}
}

func (di *RailDeviceAPI) createButton(deviceRecipe devicerecipe.RailDeviceRecipe) (button Inputer) {
	input, _ := di.boardsIOAPI.GetInputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNrPrim)
	button = raildevices.NewButton(input, deviceRecipe.Name)
	return
}

func (di *RailDeviceAPI) createToggleButton(deviceRecipe devicerecipe.RailDeviceRecipe) (toggleButton Inputer) {
	input, _ := di.boardsIOAPI.GetInputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNrPrim)
	toggleButton = raildevices.NewToggleButton(input, deviceRecipe.Name)
	return
}

func (di *RailDeviceAPI) createLamp(deviceRecipe devicerecipe.RailDeviceRecipe) (rd *runableDevice) {
	co := raildevices.NewCommonOutput(deviceRecipe.Name, getTiming(deviceRecipe))
	output, _ := di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNrPrim)
	lamp := raildevices.NewLamp(co, output)
	rd = newRunableDevice(lamp)
	return
}

func (di *RailDeviceAPI) createTwoLightSignal(deviceRecipe devicerecipe.RailDeviceRecipe) (rd *runableDevice) {
	outputPass, _ := di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNrPrim)
	outputStop, _ := di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNrSec)
	co := raildevices.NewCommonOutput(deviceRecipe.Name, getTiming(deviceRecipe))
	signal := raildevices.NewTwoLightsSignal(co, outputPass, outputStop)
	rd = newRunableDevice(signal)
	return
}

func (di *RailDeviceAPI) createTurnout(deviceRecipe devicerecipe.RailDeviceRecipe) (rd *runableDevice) {
	outputBranch, _ := di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNrPrim)
	outputMain, _ := di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNrSec)
	co := raildevices.NewCommonOutput(deviceRecipe.Name, getTiming(deviceRecipe))
	turnout := raildevices.NewTurnout(co, outputBranch, outputMain)
	rd = newRunableDevice(turnout)
	return
}

func getKey(railDeviceName string) (railDeviceKey string) {
	railDeviceKey = strings.Replace(strings.ToLower(railDeviceName), " ", "_", -1)
	return
}

func getTiming(r devicerecipe.RailDeviceRecipe) raildevices.Timing {
	start, _ := time.ParseDuration(r.StartingDelay)
	stop, _ := time.ParseDuration(r.StartingDelay)
	return raildevices.Timing{Starting: start, Stopping: stop}
}
