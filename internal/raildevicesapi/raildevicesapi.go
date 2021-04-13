package raildevicesapi

import (
	"fmt"
	"strings"

	"github.com/gen2thomas/gobrail/internal/boardpin"
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
func (di *RailDeviceAPI) AddDevice(deviceRecipe RailDeviceRecipe) (err error) {
	railDeviceKey := getKey(deviceRecipe.Name)
	if _, ok := di.devices[railDeviceKey]; ok {
		return fmt.Errorf("Rail device '%s' (key: %s) already in use", deviceRecipe.Name, railDeviceKey)
	}
	var inDev Inputer
	var runDev *runableDevice
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

func (di *RailDeviceAPI) createButton(deviceRecipe RailDeviceRecipe) (button Inputer) {
	input, _ := di.boardsIOAPI.GetInputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNr)
	button = raildevices.NewButton(input, deviceRecipe.Name)
	return
}

func (di *RailDeviceAPI) createToggleButton(deviceRecipe RailDeviceRecipe) (toggleButton Inputer) {
	input, _ := di.boardsIOAPI.GetInputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNr)
	toggleButton = raildevices.NewToggleButton(input, deviceRecipe.Name)
	return
}

func (di *RailDeviceAPI) createLamp(deviceRecipe RailDeviceRecipe) (rd *runableDevice) {
	co := raildevices.NewCommonOutput(deviceRecipe.Name, deviceRecipe.Timing)
	output, _ := di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNr)
	lamp := raildevices.NewLamp(co, output)
	rd = newRunableDevice(lamp)
	return
}

func (di *RailDeviceAPI) createTwoLightSignal(deviceRecipe RailDeviceRecipe) (rd *runableDevice) {
	outputPass, _ := di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNr)
	outputStop, _ := di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNrSecond)
	co := raildevices.NewCommonOutput(deviceRecipe.Name, deviceRecipe.Timing)
	signal := raildevices.NewTwoLightsSignal(co, outputPass, outputStop)
	rd = newRunableDevice(signal)
	return
}

func (di *RailDeviceAPI) createTurnout(deviceRecipe RailDeviceRecipe) (rd *runableDevice) {
	outputBranch, _ := di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNr)
	outputMain, _ := di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNrSecond)
	co := raildevices.NewCommonOutput(deviceRecipe.Name, deviceRecipe.Timing)
	turnout := raildevices.NewTurnout(co, outputBranch, outputMain)
	rd = newRunableDevice(turnout)
	return
}

func getKey(railDeviceName string) (railDeviceKey string) {
	railDeviceKey = strings.Replace(strings.ToLower(railDeviceName), " ", "_", -1)
	return
}
