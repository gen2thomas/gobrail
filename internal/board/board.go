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
// + most functions ready for each type of board (at the moment only typ2)
// + structure for each chip to configure
// + structure for each io at board to configure
// + set/reset all
// + set/reset one
//
// TODO:
// - each pin is analog IO with min/max value (binary can be interpreted as val=max=1/val=min=0, negotiation)
// - search for main address of board (configmode)
// - use list of already used i2cdevice addresses to exclude from search (configmode)
// - read/write eeprom at sufficient board or adaptor for "configmode"

import (
	"fmt"

	"gobot.io/x/gobot"
)

// PinType is used to type safe the constants
type PinType uint8

const (
	// Binary is used for r/w "0","1" on GPIO
	Binary PinType = iota
	// Analog is used for r/w 0-255 on analog outputs or PWM
	Analog
	// Memory is used for r/w to EEPROM
	Memory
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

type boardPin struct {
	//id to chiplist
	chipID string
	//io port of the chip
	chipPinNr uint8
	// type of the io pin
	pinType PinType
}

// PinsMap is a map of all pins on a board
type PinsMap map[uint8]*boardPin

// Board is the configuration of a board
type Board struct {
	name  string
	chips map[string]*chip
	pins  PinsMap
}

// NewBoard creates a new board with given objects
func NewBoard(name string, chips map[string]*chip, pins PinsMap) *Board {
	return &Board{name: name, chips: chips, pins: pins}
}

// GobotDevices gets all gobot devices of the board
func (b *Board) GobotDevices() []gobot.Device {
	var allDevices gobot.Devices
	for _, chip := range b.chips {
		allDevices = append(allDevices, chip.driver)
	}
	return allDevices
}

// GetBinaryPinNumbers gets all related pins of board
func (b *Board) GetBinaryPinNumbers() map[uint8]struct{} {
	return b.getPinsOfType(Binary)
}

// GetAnalogPinNumbers gets all related pins of board
func (b *Board) GetAnalogPinNumbers() map[uint8]struct{} {
	return b.getPinsOfType(Analog)
}

// GetMemoryPinNumbers gets all related pins of board
func (b *Board) GetMemoryPinNumbers() map[uint8]struct{} {
	return b.getPinsOfType(Memory)
}

// SetValue sets the given pin of board to the given value
func (b *Board) SetValue(boardPinNr uint8, value uint8) (err error) {
	var bPin *boardPin
	if bPin, err = b.getBoardPin(boardPinNr); err != nil {
		return
	}
	switch bPin.pinType {
	case Binary:
		err = b.writeGPIO(bPin, value)
	case Memory:
		err = b.writeEEPROM(bPin, value)
	default:
		err = fmt.Errorf("Pin %d with type %v not allowed to set with value %d", boardPinNr, bPin.pinType, value)
	}
	return
}

// ReadValue reads the value of the given pin of board
func (b *Board) ReadValue(boardPinNr uint8) (value uint8, err error) {
	var bPin *boardPin
	if bPin, err = b.getBoardPin(boardPinNr); err != nil {
		return
	}
	switch bPin.pinType {
	case Binary:
		value, err = b.readGPIO(bPin)
	case Memory:
		value, err = b.readEEPROM(bPin)
	default:
		err = fmt.Errorf("Pin %d with type %v not allowed to read value", boardPinNr, bPin.pinType)
	}
	return
}

// SetAllIoPins sets all pins of type "Binary" to active
func (b *Board) SetAllIoPins() (err error) {
	for ioNr, boardPin := range b.pins {
		if boardPin.pinType != Binary {
			continue
		}
		err = b.SetValue(ioNr, 0xFF)
		if err != nil {
			return err
		}
	}
	return nil
}

// ResetAllIoPins sets all pins of type "Binary" to inactive
func (b *Board) ResetAllIoPins() error {
	for ioNr, boardPin := range b.pins {
		if boardPin.pinType != Binary {
			continue
		}
		err := b.SetValue(ioNr, 0x00)
		if err != nil {
			return err
		}
	}
	return nil
}

// ShowBoardConfig prints all information of the board
func (b *Board) ShowBoardConfig() {
	fmt.Printf("\n------ Show Board (%s) ------", b)
	fmt.Printf("\n------ Chips on board ------")
	for chipID, chip := range b.chips {
		fmt.Printf("\nChip Id: %s", chipID)
		fmt.Printf(", chip driver name: %s", chip.driver.Name())
		fmt.Printf(", chip address: %d", chip.address)
	}
	fmt.Printf("\n------ Pins on board ------")
	for pinNr, boardPin := range b.pins {
		fmt.Printf("\nBoard pin number: %d", pinNr)
		fmt.Printf(", chip %s: %d (chip Id %s)", boardPin.pinType, boardPin.chipPinNr, boardPin.chipID)
	}
	fmt.Printf("\n------ Debug done ------\n")
}

func (b *Board) String() string {
	return fmt.Sprintf("Name: %s, Chips: %d, Pins: %d", b.name, len(b.chips), len(b.pins))
}

func (pt PinType) String() string {
	switch pt {
	case Memory:
		return "EEPROM address"
	case Binary:
		return "GPIO pin"
	case Analog:
		return "Ana pin"
	default:
		return "Unknown pintype"
	}
}

func (b *Board) getBoardPin(boardPinNr uint8) (boardPin *boardPin, err error) {
	var ok bool
	if boardPin, ok = b.pins[boardPinNr]; !ok {
		err = fmt.Errorf("Pin %d not there in board %s", boardPinNr, b.name)
	}
	return
}

func (b *Board) getDriver(boardPin *boardPin) (driver DriverOperations, err error) {
	var ok bool
	var chip *chip
	if chip, ok = b.chips[boardPin.chipID]; !ok {
		err = fmt.Errorf("Driver for %s not there in board %s", boardPin.chipID, b.name)
		return
	}
	driver = chip.driver
	return
}

func (b *Board) getPinsOfType(pinType PinType) (pinNumbers map[uint8]struct{}) {
	pinNumbers = make(map[uint8]struct{})
	for pinNumber, boardPin := range b.pins {
		if boardPin.pinType == pinType {
			pinNumbers[pinNumber] = struct{}{}
		}
	}
	return pinNumbers
}
