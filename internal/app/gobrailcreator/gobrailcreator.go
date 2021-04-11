package gobrailcreator

// * gets configuration from config/reader (creation plan)
// * can add elements (by json or objects?)
// * can provide creation to run railroad

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"

	"github.com/gen2thomas/gobrail/internal/boardsapi"
	"github.com/gen2thomas/gobrail/internal/raildevices"
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
	AddDevice(raildevicesapi.RailDeviceRecipe) (err error)
	ConnectNow() (err error)
	Run()
}

const boardID = "IO_Mem_PCA9501"

var boardRecipePca9501 = boardsapi.BoardRecipe{
	Name:        boardID,
	ChipDevAddr: 0x04,
	BoardType:   boardsapi.Typ2,
}

var deviceRecipeTaste1 = raildevicesapi.RailDeviceRecipe{
	Name:       "Taste 1",
	Type:       raildevicesapi.Button,
	BoardID:    boardID,
	BoardPinNr: 4,
}

var deviceRecipeTaste2 = raildevicesapi.RailDeviceRecipe{
	Name:       "Taste 2",
	Type:       raildevicesapi.ToggleButton,
	BoardID:    boardID,
	BoardPinNr: 5,
}

var deviceRecipeSignalRot = raildevicesapi.RailDeviceRecipe{
	Name:       "Signal rot",
	Type:       raildevicesapi.Lamp,
	BoardID:    boardID,
	BoardPinNr: 0,
	Timing:     raildevices.Timing{Stopping: 50 * time.Millisecond},
	Connect:    "Taste 1",
}

var deviceRecipeSignalGruen = raildevicesapi.RailDeviceRecipe{
	Name:       "Signal grün",
	Type:       raildevicesapi.Lamp,
	BoardID:    boardID,
	BoardPinNr: 3,
	Timing:     raildevices.Timing{},
	Connect:    "Signal rot",
}

var deviceRecipeSignalRotGruen = raildevicesapi.RailDeviceRecipe{
	Name:             "Rot grün Signal",
	Type:             raildevicesapi.TwoLightsSignal,
	BoardID:          boardID,
	BoardPinNr:       2,
	BoardPinNrSecond: 1,
	Timing:           raildevices.Timing{},
	Connect:          "Taste 2",
}

var boardCfgAPI BoardsConfigAPIer

// Create will create a static device connection for run
func Create(adaptor i2c.Connector) (deviceAPI RailDevicesAPIer, err error) {
	fmt.Printf("\n------ Init Boards ------\n")
	boardsAPI := boardsapi.NewBoardsAPI(adaptor)
	boardsAPI.AddBoard(boardRecipePca9501)
	fmt.Printf("\n------ Init Inputs ------\n")
	deviceAPI = raildevicesapi.NewRailDevicesAPI(boardsAPI)
	deviceAPI.AddDevice(deviceRecipeTaste1) // button
	deviceAPI.AddDevice(deviceRecipeTaste2) // togButton
	fmt.Printf("\n------ Init Outputs ------\n")
	deviceAPI.AddDevice(deviceRecipeSignalRot)      // lampY1
	deviceAPI.AddDevice(deviceRecipeSignalGruen)    // lampY2
	deviceAPI.AddDevice(deviceRecipeSignalRotGruen) // signal
	fmt.Printf("\n------ Map inputs to outputs ------\n")
	deviceAPI.ConnectNow()
	boardCfgAPI = boardsAPI
	return
}

// GobotDevices gets all gobot devices of all boards
func GobotDevices() []gobot.Device {
	return boardCfgAPI.GobotDevices()
}
