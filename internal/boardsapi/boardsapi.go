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
// - store configuration in host EEPROM or file (maybe not necessary when recipes is working)
// - detect an device and add to boards (or generate recipe automatically)
// - support for cascades
//

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"

	"github.com/gen2thomas/gobrail/internal/board"
	"github.com/gen2thomas/gobrail/internal/boardpin"
	"github.com/gen2thomas/gobrail/internal/boardrecipe"
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

// BoardsMap is the list of already created boards
type BoardsMap map[string]Boarder

// BoardsAPI is the main object for API access
type BoardsAPI struct {
	usedPins map[string]boardpin.PinNumbers
	boards   BoardsMap
	adaptor  i2c.Connector
}

// NewBoardsAPI creates a new API access
func NewBoardsAPI(adaptor i2c.Connector) *BoardsAPI {
	return &BoardsAPI{
		usedPins: make(map[string]boardpin.PinNumbers),
		boards:   make(BoardsMap),
		adaptor:  adaptor,
	}
}

// AddBoard creates a new board using recipe and add to list
func (bi *BoardsAPI) AddBoard(boardRecipe boardrecipe.Ingredients) (err error) {
	if _, ok := bi.boards[boardRecipe.Name]; ok {
		return fmt.Errorf("Board already there '%s'", boardRecipe.Name)
	}
	var newBoard Boarder
	switch boardrecipe.TypeMap[boardRecipe.Type] {
	case boardrecipe.Type2i:
		newBoard = board.NewBoardType2i(bi.adaptor, boardRecipe.ChipDevAddr, boardRecipe.Name)
	case boardrecipe.Type2o:
		newBoard = board.NewBoardType2o(bi.adaptor, boardRecipe.ChipDevAddr, boardRecipe.Name)
	case boardrecipe.Type2io:
		newBoard = board.NewBoardType2io(bi.adaptor, boardRecipe.ChipDevAddr, boardRecipe.Name)
	default:
		return fmt.Errorf("Unknown type '%s'", boardRecipe.Type)
	}
	bi.boards[boardRecipe.Name] = newBoard
	bi.usedPins[boardRecipe.Name] = make(boardpin.PinNumbers)
	return
}

// RemoveBoard remove board from list
func (bi *BoardsAPI) RemoveBoard(boardID string) {
	delete(bi.boards, boardID)
	delete(bi.usedPins, boardID)
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
	if bi.usedPins[boardID] != nil {
		bi.usedPins[boardID][boardPinNr] = struct{}{}
		return
	}
	boardPin = nil
	err = fmt.Errorf("Used pins map not initialized for %s", boardID)
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
	if bi.usedPins[boardID] != nil {
		bi.usedPins[boardID][boardPinNr] = struct{}{}
		return
	}
	boardPin = nil
	err = fmt.Errorf("Used pins map not initialized for %s", boardID)
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

// ShowAllUsedInputs list all used inputs of all boards
func (bi *BoardsAPI) ShowAllUsedInputs() {
	fmt.Printf("------ Used Pins ------\n")
	for id := range bi.boards {
		fmt.Printf("Board '%s': %s\n", id, bi.GetUsedPins(id))
	}
}

// ShowAllConfigs prints all information of all boards
func (bi *BoardsAPI) ShowAllConfigs() {
	for id := range bi.boards {
		bi.boards[id].ShowBoardConfig()
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
