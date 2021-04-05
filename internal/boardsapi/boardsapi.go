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
// + map and unmap pins
// + calculate access to the right and boardPinNr
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
	"strings"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"

	"github.com/gen2thomas/gobrail/internal/board"
)

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

type apiPin struct {
	boardID string
	// can be also a memory address
	boardPinNr uint8
	// this is the real name in difference to the key, which is all lowercase
	railDeviceName string
}

// BoardsMap is the list of already created boards
type BoardsMap map[string]*board.Board

// APIPinsMap is the type for list of pins for API
type APIPinsMap map[string]*apiPin

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

// GetFreeAPIPins gets all not already mapped API pins
func (bi *BoardsAPI) GetFreeAPIPins(boardID string, pinType board.PinType) APIPinsMap {
	var freePins = make(APIPinsMap)
	board := bi.boards[boardID]
	for boardPinNr := range board.PinsOfType(pinType) {
		if bi.FindRailDevice(boardID, boardPinNr) == "" {
			freeKey := fmt.Sprintf("Free_%s_%03d", boardID, boardPinNr)
			freePins[freeKey] = &apiPin{boardID: boardID, boardPinNr: boardPinNr}
		}
	}
	return freePins
}

// GetMappedAPIPins gets the already mapped API pins
func (bi *BoardsAPI) GetMappedAPIPins(boardID string, pinType board.PinType) APIPinsMap {
	var mappedPins = make(APIPinsMap)
	board := bi.boards[boardID]
	pinsOfType := board.PinsOfType(pinType)
	for railDeviceKey, mappedPin := range bi.mappedPins {
		if mappedPin.boardID == boardID {
			if _, ok := pinsOfType[mappedPin.boardPinNr]; ok {
				mappedPins[railDeviceKey] = mappedPin
			}
		}
	}

	return mappedPins
}

// MapPin connects a boards pin to an API pin
func (bi *BoardsAPI) MapPin(boardID string, boardPinNr uint8, railDeviceName string) {
	alreadyMappedKey := bi.FindRailDevice(boardID, boardPinNr)
	railDeviceKey := createKey(railDeviceName)
	if alreadyMappedKey == "" {
		bi.mappedPins[railDeviceKey] = &apiPin{boardID: boardID, boardPinNr: boardPinNr, railDeviceName: railDeviceName}
	} else {
		if mappedPin, ok := bi.mappedPins[railDeviceKey]; ok {
			fmt.Printf("Rail device already mapped: '%s'\n", mappedPin)
		}
		if alreadyMappedKey != railDeviceKey {
			fmt.Printf("Pin already mapped: '%s'\n", bi.mappedPins[alreadyMappedKey])
		}
	}
}

// MapPinNextFree connect a boards pin to an (randomized) free API pin (rail device)
func (bi *BoardsAPI) MapPinNextFree(boardID string, pinType board.PinType, railDeviceName string) {
	freePins := bi.GetFreeAPIPins(boardID, pinType)
	for _, freePin := range freePins {
		bi.MapPin(boardID, freePin.boardPinNr, railDeviceName)
		break
	}
}

// ReleasePin remove the connection between boards pin and an API pin (rail device)
func (bi *BoardsAPI) ReleasePin(railDeviceName string) {
	railDeviceKey := createKey(railDeviceName)
	if _, ok := bi.mappedPins[railDeviceKey]; !ok {
		fmt.Printf("Rail device '%s' not mapped, no release needed\n", railDeviceName)
	} else {
		delete(bi.mappedPins, railDeviceKey)
	}
}

// FindRailDevice gets the mapped rail device (if mapped) to a given boards pin
func (bi *BoardsAPI) FindRailDevice(boardID string, boardPinNr uint8) (railDeviceKey string) {
	for railDeviceKey, mappedPin := range bi.mappedPins {
		if mappedPin.boardID != boardID {
			continue
		}
		if mappedPin.boardPinNr != boardPinNr {
			continue
		} else {
			return railDeviceKey
		}
	}
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

// ResetAllOutputValues sets all pins of type "Binary" to inactive for all boards
func (bi *BoardsAPI) ResetAllOutputValues() (err error) {
	for _, board := range bi.boards {
		if err := board.ResetAllIoPins(); err != nil {
			return err
		}
	}
	return
}

// SetAllOutputValues sets all pins of type "Binary" to active for all boards
func (bi *BoardsAPI) SetAllOutputValues() (err error) {
	for _, board := range bi.boards {
		if err := board.SetAllIoPins(); err != nil {
			return err
		}
	}
	return
}

// SetValue sets a value of a rail device, independent of board
func (bi *BoardsAPI) SetValue(railDeviceName string, value uint8) {
	var apiPin *apiPin
	var ok bool
	if apiPin, ok = bi.mappedPins[createKey(railDeviceName)]; !ok {
		fmt.Printf("Rail device '%s' not mapped yet (key: %s)\n", railDeviceName, createKey(railDeviceName))
		return
	}
	bi.boards[apiPin.boardID].SetValue(apiPin.boardPinNr, value)
}

// GetValue gets a value of a rail device, independent of board
func (bi *BoardsAPI) GetValue(railDeviceName string) (value uint8, err error) {
	var apiPin *apiPin
	var ok bool
	if apiPin, ok = bi.mappedPins[createKey(railDeviceName)]; !ok {
		fmt.Printf("Rail device '%s' not mapped yet (key: %s)\n", railDeviceName, createKey(railDeviceName))
		return
	}
	return bi.boards[apiPin.boardID].ReadValue(apiPin.boardPinNr)
}

func (bi *BoardsAPI) String() string {
	return fmt.Sprintf("%s\n%s", bi.boards, bi.mappedPins)
}

func (apm APIPinsMap) String() string {
	toString := ""
	for railDeviceKey, apiPin := range apm {
		toString = fmt.Sprintf("%s%s -> %s\n", toString, railDeviceKey, apiPin)
	}
	if toString == "" {
		toString = "No pins mapped yet"
	}
	return toString
}

func (ap apiPin) String() string {
	return fmt.Sprintf("Board: %s, Pin: %d, Device name: %s", ap.boardID, ap.boardPinNr, ap.railDeviceName)
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

func createKey(railDeviceName string) (railDeviceKey string) {
	railDeviceKey = strings.Replace(strings.ToLower(railDeviceName), " ", "_", -1)
	return
}
