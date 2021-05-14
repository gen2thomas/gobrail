package board

// A board is the lowest level to gobot hardware drivers.
//
// in general a pin is a connection to a usable part of the board or chip
// therefore this is also the memory (pin equals the address in memory)
//
//      Author: g2t
//  Created on: 28.03.2021
// Called from: boardsapi
// Call       : some functions from gobot-adaptor (e.g. digispark)
//
// Functions:
// + most functions ready for each type of board (at the moment only type2)
// + structure for each chip to configure
// + structure for each io at board to configure
// + set/reset all
// + set/reset one
//
// TODO:
// - search for main address of board (configmode)
// - use list of already used i2c device addresses to exclude from search (configmode)
// - read/write EEPROM at sufficient board or adaptor for "configmode"

import (
	"fmt"

	"gobot.io/x/gobot"

	"github.com/gen2thomas/gobrail/internal/boardpin"
)

// DriverOperations is an interface for interact with gobot driver for chip
type DriverOperations interface {
	gobot.Driver
	Command(string) (command func(map[string]interface{}) interface{})
}

type chip struct {
	address uint8
	driver  DriverOperations
}

// PinsMap is a map of all pins on a board, the key is the board pin number
type PinsMap map[uint8]*boardpin.Pin

// Board is the configuration of a board
type Board struct {
	name    string
	chips   map[string]*chip
	pins    PinsMap
	typeTxt string
}

// NewBoard creates a new board with given objects
func NewBoard(name string, chips map[string]*chip, pins PinsMap, typeTxt string) *Board {
	return &Board{name: name, chips: chips, pins: pins, typeTxt: typeTxt}
}

// GobotDevices gets all gobot devices of the board
func (b *Board) GobotDevices() []gobot.Device {
	var allDevices gobot.Devices
	for _, chip := range b.chips {
		allDevices = append(allDevices, chip.driver)
	}
	return allDevices
}

// GetPinNumbers gets all pins of board
func (b *Board) GetPinNumbers() (pinNumbers boardpin.PinNumbers) {
	pinNumbers = make(boardpin.PinNumbers)
	for pinNumber := range b.pins {
		pinNumbers[pinNumber] = struct{}{}
	}
	return pinNumbers
}

// GetPinNumbersOfType gets board pins of given types
func (b *Board) GetPinNumbersOfType(pinTypes ...boardpin.PinType) (pinNumbers boardpin.PinNumbers) {
	pinNumbers = make(boardpin.PinNumbers)
	for pinNumber, boardPin := range b.pins {
		if boardPin.PinTypeIsOneOf(pinTypes) {
			pinNumbers[pinNumber] = struct{}{}
		}
	}
	return pinNumbers
}

// WriteValue sets the given pin of board to the given value
func (b *Board) WriteValue(boardPinNr uint8, value uint8) (err error) {
	var bPin *boardpin.Pin
	if bPin, err = b.getBoardPin(boardPinNr); err != nil {
		return
	}
	switch bPin.PinType {
	case boardpin.Binary:
		err = b.writeGPIO(bPin, value)
	case boardpin.BinaryW:
		err = b.writeGPIO(bPin, value)
	case boardpin.NBinary:
		err = b.writeGPIO(bPin, getNegatedBinaryValue(value))
	case boardpin.NBinaryW:
		err = b.writeGPIO(bPin, getNegatedBinaryValue(value))
	case boardpin.Memory:
		err = b.writeEEPROM(bPin, value)
	case boardpin.MemoryW:
		err = b.writeEEPROM(bPin, value)
	default:
		err = fmt.Errorf("Pin %d with type %v not allowed to set with value %d", boardPinNr, bPin.PinType, value)
	}
	return
}

// ReadValue reads the value of the given pin of board
func (b *Board) ReadValue(boardPinNr uint8) (value uint8, err error) {
	var bPin *boardpin.Pin
	if bPin, err = b.getBoardPin(boardPinNr); err != nil {
		return
	}
	switch bPin.PinType {
	case boardpin.Binary:
		value, err = b.readGPIO(bPin)
	case boardpin.BinaryR:
		value, err = b.readGPIO(bPin)
	case boardpin.NBinary:
		value, err = b.readGPIO(bPin)
		value = getNegatedBinaryValue(value)
	case boardpin.NBinaryR:
		value, err = b.readGPIO(bPin)
		value = getNegatedBinaryValue(value)
	case boardpin.Memory:
		value, err = b.readEEPROM(bPin)
	case boardpin.MemoryR:
		value, err = b.readEEPROM(bPin)
	default:
		err = fmt.Errorf("Pin %d with type %v not allowed to read value", boardPinNr, bPin.PinType)
	}
	return
}

// ShowBoardConfig prints all information of the board
func (b *Board) ShowBoardConfig() {
	fmt.Printf("\n------ Show Board (%s) ------\n", b)
	fmt.Printf("------ Chips on board ------\n")
	for chipID, chip := range b.chips {
		fmt.Printf("Chip Id: %s, chip driver name: %s, chip address: %d\n", chipID, chip.driver.Name(), chip.address)
	}
	fmt.Printf("------ Pins on board ------\n")
	for pinNr, boardPin := range b.pins {
		fmt.Printf("Board pin number: %d, chip %s: %d (chip Id %s)\n", pinNr, boardPin.PinType, boardPin.ChipPinNr, boardPin.ChipID)
	}
}

func (b *Board) String() string {
	return fmt.Sprintf("Name: %s, Type: %s, Chips: %d, Pins: %d", b.name, b.typeTxt, len(b.chips), len(b.pins))
}

func (b *Board) getBoardPin(boardPinNr uint8) (boardPin *boardpin.Pin, err error) {
	var ok bool
	if boardPin, ok = b.pins[boardPinNr]; !ok {
		err = fmt.Errorf("Pin %d not there in board %s", boardPinNr, b.name)
	}
	return
}

func (b *Board) getDriver(boardPin *boardpin.Pin) (driver DriverOperations, err error) {
	var ok bool
	var chip *chip
	if chip, ok = b.chips[boardPin.ChipID]; !ok {
		err = fmt.Errorf("Driver for %s not there in board %s", boardPin.ChipID, b.name)
		return
	}
	driver = chip.driver
	return
}

func getNegatedBinaryValue(value uint8) uint8 {
	if value > 0 {
		return 0
	}
	return 1
}
