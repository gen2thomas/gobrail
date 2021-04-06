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

// GetFreeAPIBinaryPins gets all not mapped API binary pins
func (bi *BoardsAPI) GetFreeAPIBinaryPins(boardID string) (freePins APIPinsMap) {
	var board Boarder
	var ok bool
	if board, ok = bi.boards[boardID]; !ok {
		return
	}
	return bi.getFreeAPIPins(boardID, board.GetBinaryPinNumbers)
}

// GetFreeAPIAnalogPins gets all not mapped API analog pins
func (bi *BoardsAPI) GetFreeAPIAnalogPins(boardID string) (freePins APIPinsMap) {
	var board Boarder
	var ok bool
	if board, ok = bi.boards[boardID]; !ok {
		return
	}
	return bi.getFreeAPIPins(boardID, board.GetAnalogPinNumbers)
}

// GetFreeAPIMemoryPins gets all not mapped API memory pins
func (bi *BoardsAPI) GetFreeAPIMemoryPins(boardID string) (freePins APIPinsMap) {
	var board Boarder
	var ok bool
	if board, ok = bi.boards[boardID]; !ok {
		return
	}
	return bi.getFreeAPIPins(boardID, board.GetMemoryPinNumbers)
}

// GetMappedAPIBinaryPins gets the already mapped API binary pins
func (bi *BoardsAPI) GetMappedAPIBinaryPins(boardID string) (mappedPins APIPinsMap) {
	var board Boarder
	var ok bool
	if board, ok = bi.boards[boardID]; !ok {
		return
	}
	return bi.getMappedAPIPins(boardID, board.GetBinaryPinNumbers)
}

// GetMappedAPIAnalogPins gets the already mapped API analog pins
func (bi *BoardsAPI) GetMappedAPIAnalogPins(boardID string) (mappedPins APIPinsMap) {
	var board Boarder
	var ok bool
	if board, ok = bi.boards[boardID]; !ok {
		return
	}
	return bi.getMappedAPIPins(boardID, board.GetAnalogPinNumbers)
}

// GetMappedAPIMemoryPins gets the already mapped API memory pins
func (bi *BoardsAPI) GetMappedAPIMemoryPins(boardID string) (mappedPins APIPinsMap) {
	mappedPins = make(APIPinsMap)
	var board Boarder
	var ok bool
	if board, ok = bi.boards[boardID]; !ok {
		return
	}
	return bi.getMappedAPIPins(boardID, board.GetMemoryPinNumbers)
}

// MapBinaryPin connect a binary boards pin to an API pin (rail device)
func (bi *BoardsAPI) MapBinaryPin(boardID string, boardPinNr uint8, railDeviceName string) (err error) {
	return bi.mapPin(boardID, int(boardPinNr), railDeviceName, bi.GetFreeAPIBinaryPins)
}

// MapAnalogPin connect a analog boards pin to API pin (rail device)
func (bi *BoardsAPI) MapAnalogPin(boardID string, boardPinNr uint8, railDeviceName string) (err error) {
	return bi.mapPin(boardID, int(boardPinNr), railDeviceName, bi.GetFreeAPIAnalogPins)
}

// MapMemoryPin connect a memory boards pin to an API pin (rail device)
// when boardPinNr is negative a randomized board pin will be connected
func (bi *BoardsAPI) MapMemoryPin(boardID string, boardPinNrOrNegative int, railDeviceName string) (err error) {
	return bi.mapPin(boardID, boardPinNrOrNegative, railDeviceName, bi.GetFreeAPIMemoryPins)
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

func (bi *BoardsAPI) mapPin(boardID string, boardPinNrOrNegative int, railDeviceName string, f func(boardID string) APIPinsMap) (err error) {
	// raildevice mapped?
	railDeviceKey := createKey(railDeviceName)
	if mappedPin, ok := bi.mappedPins[railDeviceKey]; ok {
		return fmt.Errorf("Rail device '%s' (key: %s) already mapped: '%s'", railDeviceName, railDeviceKey, mappedPin)
	}
	// free pins?
	freePins := f(boardID)
	if len(freePins) == 0 {
		return fmt.Errorf("No free pin at '%s' to map '%s'", boardID, railDeviceName)
	}
	// get randomized free pin
	var boardPinNr uint8
	if boardPinNrOrNegative < 0 {
		for _, freePin := range freePins {
			boardPinNr = freePin.boardPinNr
			break
		}
	} else {
		boardPinNr = uint8(boardPinNrOrNegative)
	}
	// already mapped?
	alreadyMappedKey := bi.FindRailDevice(boardID, boardPinNr)
	if alreadyMappedKey != "" {
		return fmt.Errorf("Pin already mapped: '%s'", bi.mappedPins[alreadyMappedKey])
	}
	// map it
	bi.mappedPins[railDeviceKey] = &apiPin{boardID: boardID, boardPinNr: boardPinNr, railDeviceName: railDeviceName}
	return
}

func (bi *BoardsAPI) getFreeAPIPins(boardID string, f func() map[uint8]struct{}) (freePins APIPinsMap) {
	freePins = make(APIPinsMap)
	for boardPinNr := range f() {
		if bi.FindRailDevice(boardID, boardPinNr) == "" {
			freeKey := fmt.Sprintf("Free_%s_%03d", boardID, boardPinNr)
			freePins[freeKey] = &apiPin{boardID: boardID, boardPinNr: boardPinNr}
		}
	}
	return freePins
}

func (bi *BoardsAPI) getMappedAPIPins(boardID string, f func() map[uint8]struct{}) (mappedPins APIPinsMap) {
	boardPinNumbers := f()
	mappedPins = make(APIPinsMap)
	for railDeviceKey, mappedPin := range bi.mappedPins {
		if mappedPin.boardID == boardID {
			if _, ok := boardPinNumbers[mappedPin.boardPinNr]; ok {
				mappedPins[railDeviceKey] = mappedPin
			}
		}
	}

	return mappedPins
}
