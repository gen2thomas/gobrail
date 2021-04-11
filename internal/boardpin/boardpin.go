package boardpin

// A board pin describes one hardware IO connection of a board
// The board pin number is mostely different from chip pin number
// A board can contain more than one chips

import (
	"fmt"
)

// PinType is used to type safe the constants
type PinType uint8

const (
	// Binary is used for r/w "0","1" on GPIO
	Binary PinType = iota
	// BinaryR is used for read only "0","1" from GPIO
	BinaryR
	// BinaryW is used for write only "0","1" to GPIO
	BinaryW
	// NBinary is used for r/w "0","1" on negotiated GPIO
	NBinary
	// NBinaryR is used for read only "0","1" from negotiated GPIO
	NBinaryR
	// NBinaryW is used for write only "0","1" to negotiated GPIO
	NBinaryW
	// Analog is used for r/w 0-255 on analog outputs or PWM
	Analog
	// AnalogR is used for read only 0-255 from analog inputs
	AnalogR
	// AnalogW is used for write only 0-255 to analog outputs or PWM
	AnalogW
	// Memory is used for r/w to EEPROM
	Memory
	// MemoryR is used for read only from EEPROM
	MemoryR
	// MemoryW is used for write only to EEPROM
	MemoryW
)

// Pin is the description of a board pin
type Pin struct {
	ChipID    string
	ChipPinNr uint8
	PinType   PinType
	MinVal    uint8
	MaxVal    uint8
}

// Input describes an pin for reading values
type Input struct {
	BoardID string
	// can be also a memory address
	BoardPinNr uint8
	ReadValue  func() (value uint8, err error)
}

// Output describes an pin for writing values
type Output struct {
	BoardID string
	// can be also a memory address
	BoardPinNr uint8
	WriteValue func(value uint8) (err error)
}

// PinNumbers is used to store numbers, e.g. as list of free or used board pins
type PinNumbers map[uint8]struct{}

// ContainsPinType check a list of pin types contains a pin type
func ContainsPinType(pinTypes []PinType, pinTypeToSearchFor PinType) bool {
	for _, pinType := range pinTypes {
		if pinType == pinTypeToSearchFor {
			return true
		}
	}
	return false
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

func (pns PinNumbers) String() (toString string) {
	for pn := range pns {
		toString = fmt.Sprintf("%s%d, ", toString, pn)
	}
	toString = fmt.Sprintf("%s\n", toString)
	return
}
