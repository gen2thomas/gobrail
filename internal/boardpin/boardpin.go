package boardpin

// A board pin describes one hardware IO connection of a board
// The board pin number is mostly different from chip pin number
// A board can contain more than one chips

import (
	"fmt"
	"strings"
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

// PinTypeMsgMap translate pin type to a small text
var PinTypeMsgMap = map[PinType]string{
	Binary:   "Binary (GPIO pin)",
	BinaryR:  "BinaryR (GPIO pin read-only)",
	BinaryW:  "BinaryW (GPIO pin write-only)",
	NBinary:  "NBinary (negated GPIO pin)",
	NBinaryR: "NBinaryR (negated GPIO pin read-only)",
	NBinaryW: "NBinaryW (negated GPIO pin write-only)",
	Analog:   "Analog (Ana pin)",
	AnalogR:  "AnalogR (Ana pin read-only)",
	AnalogW:  "AnalogW (Ana pin write-only)",
	Memory:   "Memory (EEPROM address)",
	MemoryR:  "MemoryR (EEPROM address read-only)",
	MemoryW:  "MemoryW (EEPROM address write-only)",
}

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

// PinTypeIsOneOf returns true when the type of the pin is in the list, otherwise false
func (p Pin) PinTypeIsOneOf(pinTypes []PinType) bool {
	pinTypeToSearchFor := p.PinType
	for _, pinType := range pinTypes {
		if pinType == pinTypeToSearchFor {
			return true
		}
	}
	return false
}

func (pt PinType) String() (str string) {
	if str, ok := PinTypeMsgMap[pt]; ok {
		return str
	}
	return "Unknown pin-type"
}

func (pns PinNumbers) String() (toString string) {
	var sb strings.Builder
	for pn := range pns {
		sb.WriteString(fmt.Sprintf("%d, ", pn))
	}
	return sb.String()
}
