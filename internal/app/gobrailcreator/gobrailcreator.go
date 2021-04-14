package gobrailcreator

// * gets configuration from config/reader (creation plan)
// * can add elements (by json or objects?)
// * can provide creation to run railroad

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"

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

const boardID = "IO_Mem_PCA9501"

var boardRecipePca9501 = boardsapi.BoardRecipe{
	Name:        boardID,
	ChipDevAddr: 0x04,
	BoardType:   boardsapi.Typ2,
}

var boardCfgAPI BoardsConfigAPIer

// Create will create a static device connection for run
func Create(adaptor i2c.Connector, planFile string, deviceFiles ...string) (deviceAPI RailDevicesAPIer, err error) {
	fmt.Printf("\n------ Init APIs ------\n")
	boardsAPI := boardsapi.NewBoardsAPI(adaptor)
	deviceAPI = raildevicesapi.NewRailDevicesAPI(boardsAPI)
	boardCfgAPI = boardsAPI
	fmt.Printf("\n------ Init Boards ------\n")
	boardCfgAPI.AddBoard(boardRecipePca9501)
	fmt.Printf("\n------ Read Plan ------\n")
	deviceRecipes, err := devicerecipe.ReadPlan(planFile)
	fmt.Printf("\n------ Read and add some device recipes ------\n")
	for _, deviceFile := range deviceFiles {
		deviceRecipe, _ := devicerecipe.ReadDevice(deviceFile)
		deviceRecipes = append(deviceRecipes, deviceRecipe)
	}
	fmt.Printf("\n------ Add devices from recipe list ------\n")
	for _, deviceRecipe := range deviceRecipes {
		deviceAPI.AddDevice(deviceRecipe)
	}
	fmt.Printf("\n------ Map inputs to outputs ------\n")
	deviceAPI.ConnectNow()

	return
}

// GobotDevices gets all gobot devices of all boards
func GobotDevices() []gobot.Device {
	return boardCfgAPI.GobotDevices()
}
