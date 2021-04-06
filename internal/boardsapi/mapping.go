package boardsapi

//
// The maping part of BoardsAPI is used to handle connections between
// boards pins (e.g. "board 1", pin 3) and rail device pins (e.g. "Signal 1 red")
//
// Functions:
// + map and unmap pins
// + helper functions to work with mapped pins
//
// TODO:
// - support for cascades
//

import (
	"fmt"
	"strings"

	"github.com/gen2thomas/gobrail/internal/board"
)

type apiPin struct {
	boardID string
	// can be also a memory address
	boardPinNr uint8
	// this is the real name in difference to the key, which is all lowercase
	railDeviceName string
}

// APIPinsMap is the type for list of pins for API
type APIPinsMap map[string]*apiPin

// FindRailDevice gets the key of mapped rail device to a given boards pin, if not mapped an empty string
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

// GetFreeAPIPins gets all not already mapped API pins
func (bi *BoardsAPI) GetFreeAPIPins(boardID string, pinType board.PinType) (freePins APIPinsMap) {
	freePins = make(APIPinsMap)
	var board *board.Board
	var ok bool
	if board, ok = bi.boards[boardID]; !ok {
		return
	}
	for boardPinNr := range board.PinsOfType(pinType) {
		if bi.FindRailDevice(boardID, boardPinNr) == "" {
			freeKey := fmt.Sprintf("Free_%s_%03d", boardID, boardPinNr)
			freePins[freeKey] = &apiPin{boardID: boardID, boardPinNr: boardPinNr}
		}
	}
	return freePins
}

// GetMappedAPIPins gets the already mapped API pins
func (bi *BoardsAPI) GetMappedAPIPins(boardID string, pinType board.PinType) (mappedPins APIPinsMap) {
	mappedPins = make(APIPinsMap)
	var board *board.Board
	var ok bool
	if board, ok = bi.boards[boardID]; !ok {
		return
	}
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
func (bi *BoardsAPI) MapPin(boardID string, boardPinNr uint8, railDeviceName string) (err error) {
	railDeviceKey := createKey(railDeviceName)
	if mappedPin, ok := bi.mappedPins[railDeviceKey]; ok {
		return fmt.Errorf("Rail device '%s' (key: %s) already mapped: '%s'", railDeviceName, railDeviceKey, mappedPin)
	}
	alreadyMappedKey := bi.FindRailDevice(boardID, boardPinNr)
	if alreadyMappedKey != "" {
		return fmt.Errorf("Pin already mapped: '%s'", bi.mappedPins[alreadyMappedKey])
	}

	bi.mappedPins[railDeviceKey] = &apiPin{boardID: boardID, boardPinNr: boardPinNr, railDeviceName: railDeviceName}
	return
}

// MapPinNextFree connect a boards pin of the given type to an (randomized) free API pin (rail device)
func (bi *BoardsAPI) MapPinNextFree(boardID string, pinType board.PinType, railDeviceName string) (err error) {
	freePins := bi.GetFreeAPIPins(boardID, pinType)
	if len(freePins) == 0 {
		return fmt.Errorf("No free pin at '%s' for pin type '%d' to map '%s'", boardID, pinType, railDeviceName)
	}
	for _, freePin := range freePins {
		err = bi.MapPin(boardID, freePin.boardPinNr, railDeviceName)
		break
	}
	return
}

// ReleasePin remove the connection between boards pin and an API pin (rail device)
func (bi *BoardsAPI) ReleasePin(railDeviceName string) (err error) {
	railDeviceKey := createKey(railDeviceName)
	if _, ok := bi.mappedPins[railDeviceKey]; !ok {
		return fmt.Errorf("Rail device '%s' not mapped, no release needed", railDeviceName)
	}

	delete(bi.mappedPins, railDeviceKey)
	return
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

func createKey(railDeviceName string) (railDeviceKey string) {
	railDeviceKey = strings.Replace(strings.ToLower(railDeviceName), " ", "_", -1)
	return
}
