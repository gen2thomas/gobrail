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

type connection struct {
	name     string
	inversed bool
}

// RailDeviceAPI describes the API
type RailDeviceAPI struct {
	boardsIOAPI    BoardsIOAPIer
	devices        map[string]struct{}
	runableDevices map[string]*runableDevice
	inputDevices   map[string]Inputer
	connections    map[string]connection
}

// NewRailDevicesAPI creates a new instance of rail device API
func NewRailDevicesAPI(boardsIOAPI BoardsIOAPIer) *RailDeviceAPI {
	return &RailDeviceAPI{
		devices:        make(map[string]struct{}),
		boardsIOAPI:    boardsIOAPI,
		runableDevices: make(map[string]*runableDevice),
		inputDevices:   make(map[string]Inputer),
		connections:    make(map[string]connection),
	}
}

// AddDevice creates a device from recipe and add it to the list
func (di *RailDeviceAPI) AddDevice(deviceRecipe devicerecipe.Ingredients) (err error) {
	railDeviceKey := getKey(deviceRecipe.Name)
	if _, ok := di.devices[railDeviceKey]; ok {
		return fmt.Errorf("Rail device '%s' (key: %s) already in use", deviceRecipe.Name, railDeviceKey)
	}
	var inDev Inputer
	var runDev *runableDevice
	switch devicerecipe.TypeMap[deviceRecipe.Type] {
	case devicerecipe.Button:
		if inDev, err = di.createButton(deviceRecipe); err != nil {
			return
		}
	case devicerecipe.ToggleButton:
		if inDev, err = di.createToggleButton(deviceRecipe); err != nil {
			return
		}
	case devicerecipe.Lamp:
		if runDev, err = di.createLamp(deviceRecipe); err != nil {
			return
		}
	case devicerecipe.TwoLightsSignal:
		if runDev, err = di.createTwoLightSignal(deviceRecipe); err != nil {
			return
		}
	case devicerecipe.Turnout:
		if runDev, err = di.createTurnout(deviceRecipe); err != nil {
			return
		}
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
		di.connections[railDeviceKey] = connection{name: getKey(deviceRecipe.Connect), inversed: false}
	}
	di.devices[railDeviceKey] = struct{}{}
	return
}

// ConnectNow create all connections
func (di *RailDeviceAPI) ConnectNow() (err error) {
	for runningDevKey, runableDevice := range di.runableDevices {
		var conn connection
		var ok bool
		if conn, ok = di.connections[runningDevKey]; !ok {
			continue
		}
		var conDev Inputer
		if conDev, ok = di.runableDevices[conn.name]; !ok {
			conDev = di.inputDevices[conn.name]
		}
		if conDev != nil {
			if err := runableDevice.Connect(conDev, conn.inversed); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Device with key '%s' to connect with '%s' not found", conn.name, runableDevice.RailDeviceName())
		}
	}
	return
}

// Run calls the run functions of all runnable devices
func (di *RailDeviceAPI) Run() (err error) {
	for _, runableDevice := range di.runableDevices {
		if err = runableDevice.Run(); err != nil {
			return err
		}
	}
	return
}

func (di *RailDeviceAPI) createButton(deviceRecipe devicerecipe.Ingredients) (button Inputer, err error) {
	var input *boardpin.Input
	if input, err = di.boardsIOAPI.GetInputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNrPrim); err != nil {
		return
	}
	button = raildevices.NewButton(input, deviceRecipe.Name)
	return
}

func (di *RailDeviceAPI) createToggleButton(deviceRecipe devicerecipe.Ingredients) (toggleButton Inputer, err error) {
	var input *boardpin.Input
	if input, err = di.boardsIOAPI.GetInputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNrPrim); err != nil {
		return
	}
	toggleButton = raildevices.NewToggleButton(input, deviceRecipe.Name)
	return
}

func (di *RailDeviceAPI) createLamp(deviceRecipe devicerecipe.Ingredients) (rd *runableDevice, err error) {
	var output *boardpin.Output
	if output, err = di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNrPrim); err != nil {
		return
	}
	co := raildevices.NewCommonOutput(deviceRecipe.Name, getTiming(deviceRecipe))
	lamp := raildevices.NewLamp(co, output)
	rd = newRunableDevice(lamp)
	return
}

func (di *RailDeviceAPI) createTwoLightSignal(deviceRecipe devicerecipe.Ingredients) (rd *runableDevice, err error) {
	var outputPass *boardpin.Output
	if outputPass, err = di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNrPrim); err != nil {
		return
	}
	var outputStop *boardpin.Output
	if outputStop, err = di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNrSec); err != nil {
		return
	}
	co := raildevices.NewCommonOutput(deviceRecipe.Name, getTiming(deviceRecipe))
	signal := raildevices.NewTwoLightsSignal(co, outputPass, outputStop)
	rd = newRunableDevice(signal)
	return
}

func (di *RailDeviceAPI) createTurnout(deviceRecipe devicerecipe.Ingredients) (rd *runableDevice, err error) {
	var outputBranch *boardpin.Output
	if outputBranch, err = di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNrPrim); err != nil {
		return
	}
	var outputMain *boardpin.Output
	if outputMain, err = di.boardsIOAPI.GetOutputPin(deviceRecipe.BoardID, deviceRecipe.BoardPinNrSec); err != nil {
		return
	}
	timing := getTiming(deviceRecipe)
	timing.Limit(time.Duration(1 * time.Second))
	co := raildevices.NewCommonOutput(deviceRecipe.Name, timing)
	turnout := raildevices.NewTurnout(co, outputBranch, outputMain)
	rd = newRunableDevice(turnout)
	return
}

func getKey(railDeviceName string) (railDeviceKey string) {
	railDeviceKey = strings.Replace(strings.ToLower(railDeviceName), " ", "_", -1)
	return
}

func getTiming(r devicerecipe.Ingredients) raildevices.Timing {
	start, _ := time.ParseDuration(r.StartingDelay)
	stop, _ := time.ParseDuration(r.StartingDelay)
	return raildevices.Timing{Starting: start, Stopping: stop}
}
