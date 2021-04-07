package boardsapi

//
// The BoardsAPI is the application interface from rail devices to boards
// Rail device pins can be mapped to boards IO's.
//
//      Author: g2t
//  Created on: 13.06.2009
//    Modified: 18.04.2013
//   in golang: 30.03.2021
// Called from: extDev-Modellbahn.cpp (outdated)
// Call       : some functions from board package
//
// Functions:
// + calculate access to the right board and boardPinNr
// + reset and set functions fore all boards
//
// TODO:
// - split generate recipes or read recipes from config
// - support for cascades
// - store configuration in host eeprom or file (maybe not necessary when recipes is working)
// - detect an device and add to boards (or generate recipe automatically)
//

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"

	"github.com/gen2thomas/gobrail/internal/board"
)

// ConfigurationOperations is an interface for interact with configuration part
type ConfigurationOperations interface {
	ShowBoardConfig()
}

// Boarder is an interface for interact with a board
type Boarder interface {
	ConfigurationOperations
	GobotDevices() []gobot.Device
	GetBinaryPinNumbers() map[uint8]struct{}
	GetAnalogPinNumbers() map[uint8]struct{}
	GetMemoryPinNumbers() map[uint8]struct{}
	ReadValue(boardPinNr uint8) (uint8, error)
	SetValue(boardPinNr uint8, value uint8) (err error)
}

type boardType uint8

const (
	// Typ2 is the board with a single PCA9501 and 4 amplified outputs
	Typ2 boardType = iota
	// TypUnknown is fo fallback
	TypUnknown
)

// BoardRecipe is a short description to create a new board
type BoardRecipe struct {
	Name        string
	ChipDevAddr uint8
	BoardType   boardType
}

// BoardsMap is the list of already created boards
type BoardsMap map[string]Boarder

// BoardsAPI is the main object for API access
type BoardsAPI struct {
	mappedPins APIPinsMap
	boards     BoardsMap
}

// NewBoardsAPI creates a new API access
func NewBoardsAPI(adaptor i2c.Connector, boardRecipes []BoardRecipe) *BoardsAPI {
	allBoards := make(BoardsMap)
	for _, boardRecipe := range boardRecipes {
		switch boardRecipe.BoardType {
		case Typ2:
			newBoard := board.NewBoardTyp2(adaptor, boardRecipe.ChipDevAddr, boardRecipe.Name)
			allBoards[boardRecipe.Name] = newBoard
		default:
			fmt.Println("Unknown type", boardRecipe.BoardType)
		}
	}

	bi := &BoardsAPI{
		mappedPins: make(APIPinsMap),
		boards:     allBoards,
	}
	return bi
}

// GobotDevices gets all gobot devices of all boards
func (bi *BoardsAPI) GobotDevices() []gobot.Device {
	var allDevices gobot.Devices
	for _, board := range bi.boards {
		allDevices = append(allDevices, board.GobotDevices()...)
	}
	return allDevices
}

// ShowAvailableBoards list all created boards
func (bi *BoardsAPI) ShowAvailableBoards() {
	fmt.Println(bi.boards)
}

// ShowConfig prints all information of a board
func (bi *BoardsAPI) ShowConfig(boardID string) {
	fmt.Printf("Board Id: %s\n", boardID)
	bi.boards[boardID].ShowBoardConfig()
}

// ShowConfigs prints all information of all boards
func (bi *BoardsAPI) ShowConfigs() {
	for id := range bi.boards {
		bi.ShowConfig(id)
	}
}

// SetValue sets a value of a rail device, independent of board
func (bi *BoardsAPI) SetValue(railDeviceName string, value uint8) (err error) {
	var apiPin *apiPin
	var ok bool
	if apiPin, ok = bi.mappedPins[createKey(railDeviceName)]; !ok {
		return fmt.Errorf("Rail device '%s' not mapped yet (key: %s)", railDeviceName, createKey(railDeviceName))
	}
	return bi.boards[apiPin.boardID].SetValue(apiPin.boardPinNr, value)
}

// GetValue gets a value of a rail device, independent of board
func (bi *BoardsAPI) GetValue(railDeviceName string) (value uint8, err error) {
	var apiPin *apiPin
	var ok bool
	if apiPin, ok = bi.mappedPins[createKey(railDeviceName)]; !ok {
		return 0, fmt.Errorf("Rail device '%s' not mapped yet (key: %s)", railDeviceName, createKey(railDeviceName))
	}
	return bi.boards[apiPin.boardID].ReadValue(apiPin.boardPinNr)
}

func (bi *BoardsAPI) String() string {
	return fmt.Sprintf("%s\n%s", bi.boards, bi.mappedPins)
}

func (bm BoardsMap) String() string {
	countBoards := len(bm)
	if countBoards == 0 {
		return "No Boards"
	}
	toString := fmt.Sprintf("Boards: %d\n", countBoards)
	for id, board := range bm {
		toString = fmt.Sprintf("%sBoard Id: %s, %s\n", toString, id, board)
	}
	return toString
}
