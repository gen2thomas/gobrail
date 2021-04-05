package board

// A board is the lowest level to gobot hardware drivers.
//
// in general a pin is a connection to a usable part of the board or chip
// therefore this is also the memory (pin equals the address in memory)
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

// ChipOperations is an interface for interact with gobot driver for chip
type ChipOperations interface {
	gobot.Driver
	Command(string) (command func(map[string]interface{}) interface{})
}

// ConfigurationOperations is an interface for interact with configuration part
type ConfigurationOperations interface {
	WriteBoardConfig() (err error)
	ReadBoardConfig() (err error)
	ShowBoardConfig()
}

// DeviceOperations is an interface for interact with underlying device
type DeviceOperations interface {
	ConfigurationOperations
	ReadValue(boardPinNr uint8) (uint8, error)
	SetValue(boardPinNr uint8, value uint8) (err error)
	SetAllIoPins() (err error)
	ResetAllIoPins() (err error)
	writeGPIO(pin uint8, val uint8) (err error)
	readGPIO(pin uint8) (val uint8, err error)
	writeEEPROM(address uint8, val uint8) (err error)
	readEEPROM(address uint8) (val uint8, err error)
}

type chip struct {
	address uint8
	device  ChipOperations
}

type boardPin struct {
	//id to chiplist
	chipID string
	//io port of the chip
	chipPin uint8
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

// Devices gets all devices of the board
func (b *Board) Devices() []gobot.Device {
	var allDevices gobot.Devices
	for _, chip := range b.chips {
		allDevices = append(allDevices, chip.device)
	}
	return allDevices
}

// PinsOfType gets all pins of board for the given type
func (b *Board) PinsOfType(pinType PinType) PinsMap {
	pins := make(PinsMap)
	for idx, boardPin := range b.pins {
		if boardPin.pinType == pinType {
			pins[idx] = boardPin
		}
	}
	return pins
}

// SetValue sets the given pin of board to the given value
func (b *Board) SetValue(boardPinNr uint8, value uint8) (err error) {
	//get actual device first (can be main or casc)
	var bPin *boardPin
	var ok bool
	if bPin, ok = b.pins[boardPinNr]; !ok {
		err = fmt.Errorf("Pin %d not there in board %s", boardPinNr, b.name)
	}
	switch bPin.pinType {
	case Binary:
		err = b.writeGPIO(bPin.chipPin, value)
	case Memory:
		err = b.writeEEPROM(bPin.chipPin, value)
	default:
		err = fmt.Errorf("Pin %d with type %v not allowed to set with value %d", boardPinNr, bPin.pinType, value)
	}
	return
}

// ReadValue reads the value of the given pin of board
func (b *Board) ReadValue(boardPinNr uint8) (value uint8, err error) {
	//get actual device first (can be main or casc)
	var bPin *boardPin
	var ok bool
	if bPin, ok = b.pins[boardPinNr]; !ok {
		err = fmt.Errorf("Pin %d not there in board %s", boardPinNr, b.name)
	}
	switch bPin.pinType {
	case Binary:
		value, err = b.readGPIO(bPin.chipPin)
	case Memory:
		value, err = b.readEEPROM(bPin.chipPin)
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
		fmt.Printf(", chip driver name: %s", chip.device.Name())
		fmt.Printf(", chip address: %d", chip.address)
	}
	fmt.Printf("\n------ Pins on board ------")
	for pinNr, boardPin := range b.pins {
		fmt.Printf("\nBoard pin number: %d", pinNr)
		fmt.Printf(", chip %s: %d (chip Id %s)", boardPin.pinType, boardPin.chipPin, boardPin.chipID)
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
