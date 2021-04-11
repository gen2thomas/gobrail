package boardsapi

//
// The BoardsAPI is the application interface to boards
//
//      Author: g2t
//  Created on: 13.06.2009
//    Modified: 18.04.2013
//   in golang: 30.03.2021
// Called from: extDev-Modellbahn.cpp (outdated)
// Call       : some functions from board package
//
// Functions:
// + get input and output pins and mark used
// + get all pin numbers of a board
// + get used pin numbers of a board
// + get available pin numbers of a board
//
// TODO:
// - release pins (remove used mark)
// - split generate recipes or read recipes from config
// - store configuration in host eeprom or file (maybe not necessary when recipes is working)
// - detect an device and add to boards (or generate recipe automatically)
// - support for cascades
//

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"

	"github.com/gen2thomas/gobrail/internal/board"
	"github.com/gen2thomas/gobrail/internal/boardpin"
)

// ConfigurationOperations is an interface for interact with configuration part
type ConfigurationOperations interface {
	ShowBoardConfig()
}

// Boarder is an interface for interact with a board
type Boarder interface {
	ConfigurationOperations
	GobotDevices() []gobot.Device
	GetPinNumbers() boardpin.PinNumbers
	ReadValue(boardPinNr uint8) (uint8, error)
	WriteValue(boardPinNr uint8, value uint8) (err error)
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
	usedPins map[string]boardpin.PinNumbers
	boards   BoardsMap
}

// NewBoardsAPI creates a new API access
func NewBoardsAPI(adaptor i2c.Connector, boardRecipes []BoardRecipe) *BoardsAPI {
	bi := &BoardsAPI{
		usedPins: make(map[string]boardpin.PinNumbers),
		boards:   make(BoardsMap),
	}
	for _, boardRecipe := range boardRecipes {
		switch boardRecipe.BoardType {
		case Typ2:
			newBoard := board.NewBoardTyp2(adaptor, boardRecipe.ChipDevAddr, boardRecipe.Name)
			bi.boards[boardRecipe.Name] = newBoard
			bi.usedPins[boardRecipe.Name] = make(boardpin.PinNumbers)
		default:
			fmt.Println("Unknown type", boardRecipe.BoardType)
		}
	}

	return bi
}

// GetFreePins gets all not used board pins
func (bi *BoardsAPI) GetFreePins(boardID string) (freePins boardpin.PinNumbers) {
	var board Boarder
	var ok bool
	if board, ok = bi.boards[boardID]; !ok {
		return
	}
	allPins := board.GetPinNumbers()
	usedPins := bi.usedPins[boardID]
	freePins = make(boardpin.PinNumbers)
	for boardPinNr := range allPins {
		if _, ok := usedPins[boardPinNr]; !ok {
			freePins[boardPinNr] = struct{}{}
		}
	}
	return freePins
}

// GetUsedPins gets all not used board pins
func (bi *BoardsAPI) GetUsedPins(boardID string) (usedPins boardpin.PinNumbers) {
	var ok bool
	if _, ok = bi.boards[boardID]; !ok {
		return
	}
	sourceUsedPins := bi.usedPins[boardID]
	usedPins = make(boardpin.PinNumbers)
	for usedPin := range sourceUsedPins {
		usedPins[usedPin] = struct{}{}
	}
	return
}

// GetInputPin gets an board pin to use for read values
func (bi *BoardsAPI) GetInputPin(boardID string, boardPinNr uint8) (boardPin *boardpin.Input, err error) {
	// already mapped
	if _, ok := bi.usedPins[boardID][boardPinNr]; ok {
		return nil, fmt.Errorf("Board Pin '%d' at '%s' already used", boardPinNr, boardID)
	}
	// create pin
	boardPin = &boardpin.Input{
		BoardID:    boardID,
		BoardPinNr: boardPinNr,
		ReadValue: func() (value uint8, err error) {
			return bi.boards[boardID].ReadValue(boardPinNr)
		},
	}
	bi.usedPins[boardID][boardPinNr] = struct{}{}
	return
}

// GetOutputPin gets an board pin to use for write values
func (bi *BoardsAPI) GetOutputPin(boardID string, boardPinNr uint8) (boardPin *boardpin.Output, err error) {
	// already mapped
	if _, ok := bi.usedPins[boardID][boardPinNr]; ok {
		return nil, fmt.Errorf("Board Pin '%d' at '%s' already used", boardPinNr, boardID)
	}
	// create pin
	boardPin = &boardpin.Output{
		BoardID:    boardID,
		BoardPinNr: boardPinNr,
		WriteValue: func(value uint8) (err error) {
			return bi.boards[boardID].WriteValue(boardPinNr, value)
		},
	}
	bi.usedPins[boardID][boardPinNr] = struct{}{}
	return
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

func (bi *BoardsAPI) String() string {
	return fmt.Sprintf("%s\n", bi.boards)
}

func (bm BoardsMap) String() (toString string) {
	countBoards := len(bm)
	if countBoards == 0 {
		return "No Boards"
	}
	toString = fmt.Sprintf("Boards: %d\n", countBoards)
	for id, board := range bm {
		toString = fmt.Sprintf("%sBoard Id: %s, %s\n", toString, id, board)
	}
	return toString
}
