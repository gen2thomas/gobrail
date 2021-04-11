package raildevicesapi

import (
	"fmt"
	"strings"

	"github.com/gen2thomas/gobrail/internal/boardpin"
	"github.com/gen2thomas/gobrail/internal/raildevices"
)

// RailDevice can run in loops
type RailDevice interface {
	raildevices.Outputer
	raildevices.Inputer
}

// BoardsIOAPIer is an interface for interact with a boards API which provides IO pins
type BoardsIOAPIer interface {
	GetInputPin(boardID string, boardPinNr uint8) (boardPin *boardpin.Input, err error)
	GetOutputPin(boardID string, boardPinNr uint8) (boardPin *boardpin.Output, err error)
}

// RailDeviceRecipe describes a recipe to creat an new rail device
type RailDeviceRecipe struct {
	Name             string
	Type             railDeviceType
	BoardID          string
	BoardPinNr       uint8
	BoardPinNrSecond uint8
	Timing           raildevices.Timing
	Connect          string
}

type railDeviceType uint8

const (
	// Button is a input device with one input
	Button railDeviceType = iota
	// ToggleButton is a input device with one input
	ToggleButton
	// Lamp is a output device with one output
	Lamp
	// TwoLightsSignal is a output device with two outputs, both outputs can't have the same state
	TwoLightsSignal
	// Turnout is a output device with two outputs
	Turnout
	// TypUnknown is fo fallback
	TypUnknown
)

// RailDeviceAPI describes the api
type RailDeviceAPI struct {
	boardsIOAPI    BoardsIOAPIer
	devices        map[string]struct{}
	runningDevices map[string]RailDevice
	inputDevices   map[string]raildevices.Inputer
	connections    map[string]string
}

// NewRailDevicesAPI creates a new instance of rail device API
func NewRailDevicesAPI(boardsIOAPI BoardsIOAPIer) *RailDeviceAPI {
	return &RailDeviceAPI{
		devices:        make(map[string]struct{}),
		boardsIOAPI:    boardsIOAPI,
		runningDevices: make(map[string]RailDevice),
		inputDevices:   make(map[string]raildevices.Inputer),
		connections:    make(map[string]string),
	}
}

// AddDevice creates a device from recipe and add it to the list
func (di *RailDeviceAPI) AddDevice(deviceRecipe RailDeviceRecipe) (err error) {
	railDeviceKey := getKey(deviceRecipe.Name)
	if _, ok := di.devices[railDeviceKey]; ok {
		return fmt.Errorf("Rail device '%s' (key: %s) already in use", deviceRecipe.Name, railDeviceKey)
	}
	var inDev raildevices.Inputer
	var runDev RailDevice
	switch deviceRecipe.Type {
	case Button:
		inDev = di.createButton(deviceRecipe)
	case ToggleButton:
		inDev = di.createToggleButton(deviceRecipe)
	case Lamp:
		runDev = di.createLamp(deviceRecipe)
	case TwoLightsSignal:
		runDev = di.createTwoLightSignal(deviceRecipe)
	case Turnout:
		runDev = di.createTurnout(deviceRecipe)
	default:
		return fmt.Errorf("Unknown type '%d'", deviceRecipe.Type)
	}
	if inDev != nil {
		di.inputDevices[railDeviceKey] = inDev
	}
	if runDev != nil {
		di.runningDevices[railDeviceKey] = runDev
	}
	if deviceRecipe.Connect != "" {
		di.connections[railDeviceKey] = getKey(deviceRecipe.Connect)
	}
	di.devices[railDeviceKey] = struct{}{}
	return
}

// ConnectNow create all connections
func (di *RailDeviceAPI) ConnectNow() (err error) {
	for runningDevKey, runningDevice := range di.runningDevices {
		var conKey string
		var ok bool
		if conKey, ok = di.connections[runningDevKey]; !ok {
			continue
		}
		var conDev raildevices.Inputer
		if conDev, ok = di.runningDevices[conKey]; !ok {
			conDev = di.inputDevices[conKey]
		}
		if conDev != nil {
			runningDevice.Connect(conDev)
		} else {
			return fmt.Errorf("Device with key '%s' to connect with '%s' not found", conKey, runningDevice.RailDeviceName())
		}
	}
	return
}

// Run calls the run functions of all devices
func (di *RailDeviceAPI) Run() {
	for _, runningDevice := range di.runningDevices {
		runningDevice.Run()
	}
}

func (di *RailDeviceAPI) createButton(deviceRecipe RailDeviceRecipe) (button raildevices.Inputer) {
	input, _ := di.boardsIOAPI.GetInputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNr)
	button = raildevices.NewButton(input, deviceRecipe.Name)
	return
}

func (di *RailDeviceAPI) createToggleButton(deviceRecipe RailDeviceRecipe) (toggleButton raildevices.Inputer) {
	input, _ := di.boardsIOAPI.GetInputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNr)
	toggleButton = raildevices.NewToggleButton(input, deviceRecipe.Name)
	return
}

func (di *RailDeviceAPI) createLamp(deviceRecipe RailDeviceRecipe) (lamp RailDevice) {
	co := raildevices.NewCommonOutput(deviceRecipe.Name, deviceRecipe.Timing, "lamp")
	output, _ := di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNr)
	lamp = raildevices.NewLamp(co, output)
	return
}

func (di *RailDeviceAPI) createTwoLightSignal(deviceRecipe RailDeviceRecipe) (signal RailDevice) {
	outputPass, _ := di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNr)
	outputStop, _ := di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNrSecond)
	co := raildevices.NewCommonOutput(deviceRecipe.Name, deviceRecipe.Timing, "two light signal")
	signal = raildevices.NewTwoLightsSignal(co, outputPass, outputStop)
	return
}

func (di *RailDeviceAPI) createTurnout(deviceRecipe RailDeviceRecipe) (turnout RailDevice) {
	outputBranch, _ := di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNr)
	outputMain, _ := di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNrSecond)
	co := raildevices.NewCommonOutput(deviceRecipe.Name, deviceRecipe.Timing, "turnout")
	turnout = raildevices.NewTurnout(co, outputBranch, outputMain)
	return
}

func getKey(railDeviceName string) (railDeviceKey string) {
	railDeviceKey = strings.Replace(strings.ToLower(railDeviceName), " ", "_", -1)
	return
}
