package board

// A board is the lowest level to gobot hardware drivers.
//
// in general a pin is a connection to a usable part of the board or chip
// therefore this is also the memory (pin equals the address in memory)

import (
	"fmt"

	"gobot.io/x/gobot"
)

type ChipType uint8

const (
	NoChip ChipType = iota
	PCA9501
	PCA9533
)

type PinType uint8

const (
	Binary PinType = iota
	Analog
	Memory
)

type PinChangeType uint8

const (
	Set PinChangeType = iota
	Reset
	Toggle
	Blink0
	Blink1
)

type ChipOperations interface {
	gobot.Driver
	WriteGPIO(pin uint8, val uint8) (err error)
	ReadGPIO(pin uint8) (val uint8, err error)
	Command(string) (command func(map[string]interface{}) interface{})
}

type ConfigurationOperations interface {
	WriteBoardConfig() (err error)
	ReadBoardConfig() (err error)
	ShowBoardConfig()
}

type BoardOperations interface {
	ConfigurationOperations
	ReadValue(boardPinNr uint8) (uint8, error)
	SetValue(boardPinNr uint8, value uint8) (err error)
	SetAllIoPins() (err error)
	ResetAllIoPins() (err error)
	writeEEPROM(address uint8, val uint8) (err error)
	readEEPROM(address uint8) (val uint8, err error)
}

type chip struct {
	chipType ChipType
	address  uint8
	device   ChipOperations
}

type boardPin struct {
	//id to chiplist
	chipId string
	//io port of the chip
	chipPin uint8
	// type of the io pin
	pinType PinType
}

type boardPinsMap map[uint8]*boardPin

type Board struct {
	name  string
	chips map[string]*chip
	pins  boardPinsMap
}

var chipEmpty = chip{
	chipType: NoChip,
	address:  0xFF,
	device:   nil,
}

// an empty board
var boardEmpty = Board{
	name:  "EmptyBoard",
	chips: map[string]*chip{"EmptyChip": &chipEmpty},
	pins: boardPinsMap{
		0: {chipId: "EmptyChip", chipPin: 0, pinType: Binary},
		1: {chipId: "EmptyChip", chipPin: 1, pinType: Binary},
		2: {chipId: "EmptyChip", chipPin: 2, pinType: Binary},
		3: {chipId: "EmptyChip", chipPin: 3, pinType: Binary},
	},
}

func (b *Board) Devices() []gobot.Device {
	var allDevices gobot.Devices
	for _, chip := range b.chips {
		allDevices = append(allDevices, chip.device)
	}
	return allDevices
}

func (b *Board) PinsOfType(pinType PinType) boardPinsMap {
	pins := make(boardPinsMap)
	for idx, boardPin := range b.pins {
		if boardPin.pinType == pinType {
			pins[idx] = boardPin
		}
	}
	return pins
}

func (b *Board) SetValue(boardPinNr uint8, value uint8) (err error) {
	//get actual device first (can be main or casc)
	var bPin *boardPin
	var ok bool
	if bPin, ok = b.pins[boardPinNr]; !ok {
		err = fmt.Errorf("Pin %d not there in board %s", boardPinNr, b.name)
	}
	chip := b.chips[bPin.chipId]
	switch bPin.pinType {
	case Binary:
		err = chip.device.WriteGPIO(bPin.chipPin, value)
	case Memory:
		err = b.writeEEPROM(bPin.chipPin, value)
	default:
		err = fmt.Errorf("Pin %d with type %v not allowed to set with value %d", boardPinNr, bPin.pinType, value)
	}
	return
}

func (b *Board) ReadValue(boardPinNr uint8) (value uint8, err error) {
	//get actual device first (can be main or casc)
	var bPin *boardPin
	var ok bool
	if bPin, ok = b.pins[boardPinNr]; !ok {
		err = fmt.Errorf("Pin %d not there in board %s", boardPinNr, b.name)
	}
	chip := b.chips[bPin.chipId]
	switch bPin.pinType {
	case Binary:
		value, err = chip.device.ReadGPIO(bPin.chipPin)
	case Memory:
		value, err = b.readEEPROM(bPin.chipPin)
	default:
		err = fmt.Errorf("Pin %d with type %v not allowed to read value", boardPinNr, bPin.pinType)
	}
	return
}

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

func (b *Board) ShowBoardConfig() {
	fmt.Printf("\n------ Show Board (%s) ------", b)
	fmt.Printf("\n------ Chips on board ------")
	for chipId, chip := range b.chips {
		fmt.Printf("\nChip Id: %s", chipId)
		fmt.Printf(", chip type: %d", chip.chipType)
		fmt.Printf(", chip address: %d", chip.address)
	}
	fmt.Printf("\n------ Pins on board ------")
	for pinNr, boardPin := range b.pins {
		fmt.Printf("\nBoard pin number: %d", pinNr)
		fmt.Printf(", chip %s: %d (chip Id %s)", boardPin.pinType, boardPin.chipPin, boardPin.chipId)
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
