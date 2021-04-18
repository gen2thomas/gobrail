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
	"gobot.io/x/gobot/platforms/digispark"
	"gobot.io/x/gobot/platforms/raspi"
	"gobot.io/x/gobot/platforms/tinkerboard"

	"github.com/gen2thomas/gobrail/internal/boardsapi"
	"github.com/gen2thomas/gobrail/internal/devicerecipe"
	"github.com/gen2thomas/gobrail/internal/raildevicesapi"
)

// BoardsConfigAPIer is an interface for interact with a boards API
type BoardsConfigAPIer interface {
	GobotDevices() []gobot.Device
	AddBoard(boardRecipe boardsapi.BoardRecipe) (err error)
	RemoveBoard(boardID string)
}

// RailDevicesAPIer is an interface to interact with rail devices
type RailDevicesAPIer interface {
	AddDevice(devicerecipe.RailDeviceRecipe) (err error)
	ConnectNow() (err error)
	Run()
}

// RailRunner is an interface to poll the rail
type RailRunner interface {
	Run()
}

type i2cAdaptor interface {
	i2c.Connector
	gobot.Connection
}

// AdaptorType represents the supported adaptors
type AdaptorType uint8

const (
	digisparkType AdaptorType = iota
	raspiType
	tinkerboardType
	unknownType
)

var adaptorTypeToStringMap = map[AdaptorType]string{digisparkType: "digispark", raspiType: "raspi", tinkerboardType: "tinkerboard", unknownType: "typUnknown"}
var adaptorStringToTypeMap = map[string]AdaptorType{"digispark": digisparkType, "raspi": raspiType, "tinkerboard": tinkerboardType, "unknown": unknownType}

const boardID = "IO_Mem_PCA9501"

var boardRecipePca9501 = boardsapi.BoardRecipe{
	Name:        boardID,
	ChipDevAddr: 0x04,
	BoardType:   boardsapi.Typ2,
}

var lastGobot *gobot.Robot

// Create will create a static device connection for run
// before creating, the old gobot robot will be stopped
// after creating the devices, a new gobot robot will be created and started
func Create(daemonMode bool, name string, adaptorType AdaptorType, planFile string, deviceFiles ...string) (runner RailRunner, err error) {
	if err = Stop(); err != nil {
		return
	}

	fmt.Printf("\n------ Init gobot adaptor (%s) ------\n", adaptorType)
	var adaptor i2cAdaptor
	if adaptor, err = createAdaptor(adaptorType); err != nil {
		return
	}
	fmt.Printf("\n------ Init APIs ------\n")
	boardsAPI := boardsapi.NewBoardsAPI(adaptor)
	deviceAPI := raildevicesapi.NewRailDevicesAPI(boardsAPI)
	fmt.Printf("\n------ Init Boards ------\n")
	boardsAPI.AddBoard(boardRecipePca9501)
	fmt.Printf("\n------ Read Plan (%s) ------\n", planFile)
	var deviceRecipes []devicerecipe.RailDeviceRecipe
	if deviceRecipes, err = devicerecipe.ReadPlan(planFile); err != nil {
		return
	}
	fmt.Printf("\n------ Read and add some device recipes ------\n")
	for _, deviceFile := range deviceFiles {
		var deviceRecipe devicerecipe.RailDeviceRecipe
		if deviceRecipe, err = devicerecipe.ReadDevice(deviceFile); err != nil {
			return
		}
		deviceRecipes = append(deviceRecipes, deviceRecipe)
	}
	fmt.Printf("\n------ Add devices from recipe list ------\n")
	for _, deviceRecipe := range deviceRecipes {
		deviceAPI.AddDevice(deviceRecipe)
	}
	fmt.Printf("\n------ Map inputs to outputs ------\n")
	deviceAPI.ConnectNow()

	if daemonMode {
		// cyclic call of "Run()" is done by daemon program
		lastGobot = gobot.NewRobot(name,
			[]gobot.Connection{adaptor},
			boardsAPI.GobotDevices(),
		)
		// very important for daemon mode
		lastGobot.AutoRun = false
	} else {
		work := func() {
			gobot.Every(50*time.Millisecond, func() {
				deviceAPI.Run()
			})
		}

		lastGobot = gobot.NewRobot(name,
			[]gobot.Connection{adaptor},
			boardsAPI.GobotDevices(),
			work,
		)
	}

	if err = lastGobot.Start(); err != nil {
		return
	}

	return deviceAPI, nil
}

// Stop stops the gobot robot, when available
func Stop() (err error) {
	if lastGobot != nil {
		if lastGobot.Running() {
			fmt.Printf("\n------ Stop gobot (%s) ------\n", lastGobot.Name)
			err = lastGobot.Stop()
		}
		lastGobot = nil
	}
	return
}

// ParseAdaptorType try get adaptor type from string
func ParseAdaptorType(adaptorString string) (a AdaptorType, err error) {
	var ok bool
	if a, ok = adaptorStringToTypeMap[strings.ToLower(adaptorString)]; !ok {
		err = fmt.Errorf("Unknown adaptor %s", adaptorString)
	}
	return
}

func (a AdaptorType) String() string {
	return adaptorTypeToStringMap[a]
}

func createAdaptor(adaptorType AdaptorType) (adaptor i2cAdaptor, err error) {
	switch adaptorType {
	case digisparkType:
		adaptor = digispark.NewAdaptor()
	case raspiType:
		adaptor = raspi.NewAdaptor()
	case tinkerboardType:
		adaptor = tinkerboard.NewAdaptor()
	default:
		err = fmt.Errorf("Unknown type '%d'", adaptorType)
	}
	return
}
